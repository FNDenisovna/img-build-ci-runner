package periodically_runner

import (
	"fmt"
	config "img-build-ci-runner/internal/config/viper"
	wf_runner "img-build-ci-runner/internal/integration/wf_runner"
	model "img-build-ci-runner/internal/model"
	"log"
	"strings"
	"time"
)

type Service struct {
	cfg *config.Config
}

func New(c *config.Config) *Service {
	return &Service{
		cfg: c,
	}
}

func (s *Service) Run(simulate bool, closing chan bool) error {
	//Get list images-packages for checking from config source
	//key - image, value - list of packs
	orgGroups := strings.Split(s.cfg.GetString(config.PeriodImgGroupCfgKey), " ")
	branches := strings.Split(s.cfg.GetString(config.BranchesCfgKey), " ")
	wfUrl := s.cfg.GetString(config.WfUrlCfgKey)
	wfOrgRepo := s.cfg.GetString(config.WfOrgRepoCfgKey)
	wfRefRepo := s.cfg.GetString(config.WfRefRepoCfgKey)
	wfName := s.cfg.GetString(config.WfGroupNameCfgKey)
	wfToken := s.cfg.GetString(config.WfTokenCfgKey)

	//foreach branch check
	for _, b := range branches {
		select {
		case <-closing:
			log.Println("Finish base images worker by closing chanell singal")
			return nil
		default:
			// foreach org send build tag
			// to build full group
			for _, o := range orgGroups {
				//generate run building image
				data := &model.WfInputDataGroup{
					Inputs: model.WfInputsGroup{
						Group:  o,
						Branch: b,
					},
					Ref: fmt.Sprintf("refs/heads/%s", wfRefRepo),
				}
				wf_runner.RunBuildImage(data, wfUrl, wfName, wfOrgRepo, wfToken)

				time.Sleep(time.Minute * 90)
			}
		}
	}

	return nil
}
