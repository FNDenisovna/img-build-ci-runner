package vers_checker_runner

import (
	"fmt"
	"img-build-ci-runner/internal/compare"
	img_info_getter "img-build-ci-runner/internal/img_info_getter/git_getter"
	model "img-build-ci-runner/internal/model"
	"log"
	"time"
)

type Service struct {
	altapi        AltApi
	gitea         GitApi
	db            Db
	branches      []string
	imgInfoGetter ImgInfoGetter
	token         string
}

type ImgInfoGetter interface {
	GetImgPkgMap() map[string][]string
	Update(giturl string, imgGroup string)
}

type Db interface {
	GetPackages(branch string, limit int) ([]model.SqlPack, error)
	InsertPackage(pack *model.SqlPack) (int, error)
}

type AltApi interface {
	GetSitePackInfo(name, branch string) (packinfo model.SiteVersion, err error)
	GetPackInfo(name, branch string) (packinfo model.SiteVersion, err error)
}

type GitApi interface {
	RunBuildImage(tag *model.GiteaTag, token string) error
}

func New(altapi AltApi, gitea GitApi, db Db, branches []string, imgPkgGetterSouce string, imgGroup string, token string) *Service {
	imgInfoGetter := img_info_getter.New(imgPkgGetterSouce, imgGroup)
	return &Service{
		altapi:        altapi,
		gitea:         gitea,
		db:            db,
		branches:      branches,
		imgInfoGetter: imgInfoGetter,
		token:         token,
	}
}

func (s *Service) Update(branches []string, imgPkgGetterSouce string, imgGroup string, token string) {
	s.branches = branches
	s.imgInfoGetter.Update(imgPkgGetterSouce, imgGroup)
}

// Check versions of packages in images
// If wersion is higther than in local memory
// Run images building by sending tag to workflow repo
func (s *Service) Run(simulate bool) error {
	//Get list images-packages for checking from config source
	//key - image, value - list of packs
	checklist := s.GetImgPkgMap()
	branches := s.branches

	//foreach branch check
	for _, b := range branches {
		packsDbMap, err := s.GetImgPkgMapDb(b)
		if err != nil {
			err = fmt.Errorf("Can't get packages list from db for branch %s. Error: %w\n", b, err)
			return err
		}

		//foreach image get current version pack from site
		for im, packs := range checklist {
			mainPack := packs[0]
			dbInfo, checked := packsDbMap[mainPack]

			curPackInfo, err := s.altapi.GetSitePackInfo(mainPack, b)
			if err != nil {
				curPackInfo, err = s.altapi.GetPackInfo(mainPack, b)
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

			if !simulate {
				//generate message to email
				//TODO

				//generate run building image
				tag := &model.GiteaTag{
					Image:   im,
					Branch:  b,
					Version: curPackInfo.Version,
				}
				s.gitea.RunBuildImage(tag, "")
			}

			time.Sleep(time.Second * 15)
		}
	}

	return nil
}

func (s *Service) GetImgPkgMap() map[string][]string {
	return s.imgInfoGetter.GetImgPkgMap()
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
