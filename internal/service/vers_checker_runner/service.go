package vers_checker_runner

import (
	"encoding/json"
	"fmt"
	"img-build-ci-runner/internal/compare"
	config "img-build-ci-runner/internal/config/viper"
	"img-build-ci-runner/internal/img_info_getter/git_getter"
	alt_api "img-build-ci-runner/internal/integration/alt_api"
	wf_runner "img-build-ci-runner/internal/integration/wf_runner"
	model "img-build-ci-runner/internal/model"
	"log"
	"strings"
	"time"
)

type Service struct {
	cfg *config.Config
	db  Db
}

type ImgInfoGetter interface {
	GetImgPkgMap() map[string][]string
}

type Db interface {
	GetPackages(branch string, limit int) ([]model.SqlPack, error)
	InsertPackage(pack *model.SqlPack) (int, error)
}

type WfApi interface {
	RunBuildImage(inputData *model.WfInputDataImages) error
}

func New(db Db, c *config.Config) *Service {
	return &Service{
		cfg: c,
		db:  db,
	}
}

// Check versions of packages in images
// If wersion is higther than in local memory
// Run images building by sending tag to workflow repo
func (s *Service) Run(simulate bool, closing chan bool) error {
	//Get branches and group images (from images templates repo)
	imgGroups := strings.Split(s.cfg.GetString(config.VersCheckImgGroupCfgKey), " ")
	branches := strings.Split(s.cfg.GetString(config.BranchesCfgKey), " ")
	altApiUrl := s.cfg.GetString(config.AltApiUrlCfgKey)
	wfUrl := s.cfg.GetString(config.WfUrlCfgKey)
	wfOrgRepo := s.cfg.GetString(config.WfOrgRepoCfgKey)
	wfRefRepo := s.cfg.GetString(config.WfRefRepoCfgKey)
	wfName := s.cfg.GetString(config.WfImagesNameCfgKey)
	wfToken := s.cfg.GetString(config.WfTokenCfgKey)

	//Get list images-packages for checking from config source
	//key - image, value - list of packs
	for ig, g := range imgGroups {
		checklist := s.GetImgPkgMap(g)

		//foreach branch check
		for ib, b := range branches {
			select {
			case <-closing:
				log.Println("Finish packages version checker worker by closing chanell singal")
				return nil
			default:
				packsDbMap, err := s.GetImgPkgMapDb(b)
				if err != nil {
					err = fmt.Errorf("Can't get packages list from db for branch %s. Error: %w\n", b, err)
					return err
				}

				data := &model.WfInputDataImages{
					Ref: fmt.Sprintf("refs/heads/%s", wfRefRepo),
					Inputs: model.WfInputsImages{
						Branch: b,
					},
				}

				//foreach image get current version pack from site
				for im, packs := range checklist {
					mainPack := packs[0]
					dbInfo, checked := packsDbMap[mainPack]

					curPackInfo, err := alt_api.GetSitePackInfo(altApiUrl, mainPack, b)
					if err != nil {
						log.Printf("Can't get package info from basealt site tasks info. Package: %s, Branch: %s, Error: %v\n", mainPack, b, err)
						log.Printf("Try again from basealt site total package info\n")
						curPackInfo, err = alt_api.GetPackInfo(altApiUrl, mainPack, b)
						if err != nil {
							log.Printf("Can't get package info from basealt site. Package: %s, Branch: %s, Error: %v\n", mainPack, b, err)
							continue
						}
					}

					// packege exists in db
					if checked && dbInfo.Version != "" {
						//compare
						dbVer := fmt.Sprintf("%d:%s-%s", dbInfo.Epoch, dbInfo.Version, dbInfo.Release)
						curVer := fmt.Sprintf("%d:%s-%s", 0, curPackInfo.Version, curPackInfo.Release)
						if compRes, _ := compare.Compare(curVer, dbVer); compRes <= 0 {
							continue
						}
					}

					//update version in db
					dbInfo.Name = mainPack
					dbInfo.Version = curPackInfo.Version
					dbInfo.Release = curPackInfo.Release
					dbInfo.Epoch = 0
					dbInfo.Changed = time.Now()
					dbInfo.Branch = b

					s.db.InsertPackage(&dbInfo)
					log.Printf("Insert to db: %v\n", dbInfo)

					data.Inputs.Images = append(data.Inputs.Images, model.WfInputsImagesInfo{
						Name:    fmt.Sprintf("%s/%s", g, im),
						Version: curPackInfo.Version,
					})
					time.Sleep(time.Second * 5)
				}

				imagesBytes, err := json.Marshal(data.Inputs.Images)
				if err != nil {
					log.Printf("Can't marshal data.Inputs.Images to string. Error: %v\n", err)
					continue
				}
				//imagesStr := strings.ReplaceAll(string(imagesBytes), "\"name\"", "\\\"name\\\"")
				//data.Inputs.ImagesStr = strings.ReplaceAll(imagesStr, "\"version\"", "\\\"version\\\"")
				data.Inputs.ImagesStr = string(imagesBytes)

				//foreach branch run building workflow
				if !simulate {
					//generate message to email
					//TODO

					_ = wf_runner.RunBuildImage(data, wfUrl, wfName, wfOrgRepo, wfToken)
				}

				if ib >= len(branches)-1 {
					continue
				} else {
					time.Sleep(time.Minute * 40)
				}
			}

			// time delay between running workflow and new
			if ig >= len(imgGroups)-1 {
				continue
			} else {
				time.Sleep(time.Minute * 40)
			}
		}
	}

	return nil
}

func (s *Service) GetImgPkgMap(imgGroup string) map[string][]string {
	imgInfoGetter := git_getter.New(s.cfg.GetString(config.ImgPkgGetterSouceCfgKey), imgGroup)
	return imgInfoGetter.GetImgPkgMap()
}

func (s *Service) GetImgPkgMapDb(b string) (map[string]model.SqlPack, error) {
	//Get packages versions from db
	packsDb, err := s.db.GetPackages(b, 0)
	if err != nil {
		return nil, err
	}

	// key - pack, value - pack model from db
	res := make(map[string]model.SqlPack)
	for _, p := range packsDb {
		if _, ok := res[p.Name]; !ok {
			res[p.Name] = p
		}
	}

	return res, nil
}
