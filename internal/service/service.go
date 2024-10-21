package service

import (
	"altpack-vers-checker/internal/compare"
	model "altpack-vers-checker/internal/model"
	"fmt"
	"log"
	"time"
)

type Service struct {
	altapi AltApi
	cfg    Cfg
	db     Db
	gitea  GitApi
}

type Db interface {
	GetPackages(branch string, limit int) ([]model.SqlPack, error)
	InsertPackage(pack *model.SqlPack) (int, error)
}

type Cfg interface {
	GetImgPkgList() ([]model.ImgPkg, error)
	GetBranches() []string
}

type AltApi interface {
	GetSitePackInfo(name, branch string) (packinfo model.SiteVersion, err error)
}

type GitApi interface {
	CreateTag(tag *model.GiteaTag, token string) error
}

func New(altapi AltApi, gitea GitApi, cfg Cfg, db Db) *Service {
	return &Service{
		altapi: altapi,
		gitea:  gitea,
		cfg:    cfg,
		db:     db,
	}
}

func (s *Service) CheckPackagesVersion() error {
	//Get list images-packages for checking from config source
	checklist, err := s.cfg.GetImgPkgList()
	if err != nil {
		return err
	}

	branches := s.cfg.GetBranches()
	//foreach branch check
	for _, b := range branches {
		//Get packages versions from db
		packsDb, err := s.db.GetPackages(b, 0)
		if err != nil {
			fmt.Errorf("Can't get packages list from db for branch %s. Error: %v\n", b, err)
		}
		//make map
		//key - image
		//value - package from local dbPacks
		packsDbMap := make(map[string]model.SqlPack, 0)
		for _, im := range checklist {
			for _, pack := range packsDb {
				if im.Package == pack.Name {
					packsDbMap[im.Image] = pack
				}
			}

			if _, ok := packsDbMap[im.Image]; !ok {
				packsDbMap[im.Image] = model.SqlPack{
					Name: im.Package,
				}
			}
		}
		//foreach image get current version pack from site
		for k, v := range packsDbMap {
			curPackInfo, err := s.altapi.GetSitePackInfo(v.Name, b)
			if err != nil {
				log.Printf("Can't get package info from basealt site. Package: %s, Branch: %s, Error: %v\n", v.Name, b, err)
			}

			if v.Version != "" {
				//compare
				dbVer := fmt.Sprintf("%d:%s-%s", v.Epoch, v.Version, v.Release)
				curVer := fmt.Sprintf("%d:%s-%s", 0, curPackInfo.Version, curPackInfo.Release)
				if compRes, _ := compare.Compare(curVer, dbVer); compRes <= 0 {
					continue
				}
			}
			//generate message to email
			//generate run building image
			tag := &model.GiteaTag{
				Image: k,
			}
			s.gitea.CreateTag(tag, "")

			//update version in db
			v.Version = curPackInfo.Version
			v.Release = curPackInfo.Release
			v.Epoch = 0
			v.Changed = time.Now()
			v.Branch = b

			s.db.InsertPackage(&v)
		}
	}

	return nil
}
