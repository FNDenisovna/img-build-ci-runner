package deps_cron_runner

import (
	img_info_getter "img-build-ci-runner/internal/img_info_getter/git_getter"
	model "img-build-ci-runner/internal/model"
	"time"
)

type ImgInfoGetter interface {
	GetImgGroup() []string
	Update(giturl string, imgGroup string)
	GetImgWithDeps() map[string]string
}

type Service struct {
	gitea         GitApi
	branches      []string
	imgInfoGetter ImgInfoGetter
	token         string
}

type GitApi interface {
	RunBuildImage(tag *model.GiteaTag, token string) error
}

func New(gitea GitApi, branches []string, imgPkgGetterSouce string, imgGroup string, token string) *Service {
	imgInfoGetter := img_info_getter.New(imgPkgGetterSouce, imgGroup)
	return &Service{
		gitea:         gitea,
		branches:      branches,
		imgInfoGetter: imgInfoGetter,
		token:         token,
	}
}

func (s *Service) Update(branches []string, imgPkgGetterSouce string, imgGroup string, token string) {
	s.branches = branches
	s.imgInfoGetter.Update(imgPkgGetterSouce, imgGroup)
}

func (s *Service) RunSeparate(simulate bool) error {
	//TODO
	return nil
}

func (s *Service) Run(simulate bool) error {
	//Get list images-packages for checking from config source
	//key - image, value - list of packs
	orgGroups := s.imgInfoGetter.GetImgGroup()
	branches := s.branches

	//foreach branch check
	for _, b := range branches {
		// foreach org send build tag
		// to build full group
		for _, o := range orgGroups {
			//generate run building image
			tag := &model.GiteaTag{
				Branch: b,
				Org:    o,
			}
			s.gitea.RunBuildImage(tag, "")

			time.Sleep(time.Second * 15)
		}
	}

	return nil
}

func (s *Service) GetImgGroup() []string {
	return s.imgInfoGetter.GetImgGroup()
}
