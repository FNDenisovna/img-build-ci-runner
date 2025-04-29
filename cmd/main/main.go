package main

import (
	"img-build-ci-runner/internal/config"
	"img-build-ci-runner/internal/storage"

	"img-build-ci-runner/internal/integration/altapi"
	"img-build-ci-runner/internal/integration/gitea"
	deps_service "img-build-ci-runner/internal/service/deps_cron_runner"
	vers_service "img-build-ci-runner/internal/service/vers_checker_runner"
	sql "img-build-ci-runner/internal/storage/sqllite"
	"log"

	cron "gopkg.in/robfig/cron.v2"
)

func main() {
	cfg := config.New()
	aapi := altapi.New(cfg.GetSettings("AltSiteUrl"))
	gapi := gitea.New(cfg.GetSettings("GiteaWfUrl"))
	versCronExp := cfg.GetSettings("VersCronGroupExp")
	depsCronExp := cfg.GetSettings("DepsCronGroupExp")

	dbDriver, err := sql.New()
	if err != nil {
		log.Fatalf("Can't create/open db. Error: %v\n", err)
		return
	}

	db := storage.New(dbDriver)
	defer func() {
		db.Close()
	}()

	versSrv := vers_service.New(aapi, gapi, db, cfg.GetBranches(), cfg.GetSettings("GiteaRepoUrl"), cfg.GetSettings("VersCheckImgGroup"), cfg.GetSettings("GiteaToken"))
	depsSrv := deps_service.New(gapi, cfg.GetBranches(), cfg.GetSettings("GiteaRepoUrl"), cfg.GetSettings("DepsCronImgGroup"), cfg.GetSettings("GiteaToken"))
	//Run first getting versions of packageList
	if err = versSrv.Run(true); err != nil {
		log.Fatalln(err)
		return
	}

	errChan := make(chan error)

	c := cron.New()
	// add versions chercher runner
	c.AddFunc(versCronExp, func() {
		if err := cfg.UpdateSettings(); err != nil {
			errChan <- err
		}
		aapi.Update(cfg.GetSettings("AltSiteUrl"))
		gapi.Update(cfg.GetSettings("GiteaWfUrl"))
		versSrv.Update(cfg.GetBranches(), cfg.GetSettings("GiteaRepoUrl"), cfg.GetSettings("VersCheckImgGroup"), cfg.GetSettings("GiteaToken"))
		if err = versSrv.Run(true); err != nil {
			errChan <- err
		}
	})
	// add dependensy runner
	c.AddFunc(depsCronExp, func() {
		if err := cfg.UpdateSettings(); err != nil {
			errChan <- err
		}
		aapi.Update(cfg.GetSettings("AltSiteUrl"))
		gapi.Update(cfg.GetSettings("GiteaWfUrl"))
		depsSrv.Update(cfg.GetBranches(), cfg.GetSettings("GiteaRepoUrl"), cfg.GetSettings("DepsCronImgGroup"), cfg.GetSettings("GiteaToken"))
		if err = depsSrv.Run(false); err != nil {
			errChan <- err
		}
	})
	c.Start()
	defer c.Stop()

	select {
	case err := <-errChan:
		log.Fatalln(err)
		return
	}
}
