package main

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	config "img-build-ci-runner/internal/config/viper"
	period_service "img-build-ci-runner/internal/service/periodically_runner"
	vers_service "img-build-ci-runner/internal/service/vers_checker_runner"
	"img-build-ci-runner/internal/storage"
	sql "img-build-ci-runner/internal/storage/sqllite"

	cron "gopkg.in/robfig/cron.v2"
)

func main() {
	cfg := config.New()

	versCronExp := cfg.GetString("vers_check_img_group")
	periodCronExp := cfg.GetString("period_cron_omg_group")
	storagePath := cfg.GetString("storage_path")

	dbDriver, err := sql.New(storagePath)
	if err != nil {
		log.Fatalf("Can't create/open db. Error: %v\n", err)
		return
	}

	db := storage.New(dbDriver)
	defer func() {
		db.Close()
	}()

	versSrv := vers_service.New(db, cfg)
	periodSrv := period_service.New(cfg)

	exit_chan := make(chan os.Signal, 1)
	signal.Notify(exit_chan, os.Interrupt, syscall.SIGTERM)
	errChan := make(chan error) // getting errors from childs channel
	closing := make(chan bool)  // closing signal channel for childs
	success := make(chan bool)

	go func() {
		if err = versSrv.Run(true, closing); err != nil {
			errChan <- err
		}
	}()

	return

	//Run first getting versions of packageList
	go func() {
		if err = versSrv.Run(true, closing); err != nil {
			log.Println(err)
			errChan <- err
			return
		}
		success <- true
	}()

	select {
	case err = <-errChan:
		log.Println("Exit by error: ", err)
		close(closing)
		close(errChan)
		return
	case <-exit_chan:
		log.Println("Exit by os.Interrupt syscall.SIGTERM")
		// initiate graceful shutdown
		close(closing)
		close(errChan)
		return
	case <-success:
		break
	}

	wg := &sync.WaitGroup{}
	c := cron.New()
	// add versions chercher runner

	c.AddFunc(versCronExp, func() {
		wg.Add(1)
		defer wg.Done()

		if err = versSrv.Run(true, closing); err != nil {
			errChan <- err
		}
	})
	// add dependensy runner
	c.AddFunc(periodCronExp, func() {
		wg.Add(1)
		defer wg.Done()

		if err = periodSrv.Run(true, closing); err != nil {
			errChan <- err
		}
	})
	c.Start()
	defer c.Stop()

	select {
	case err = <-errChan:
		log.Print(err)
		close(closing)
		wg.Wait() //wait childs
		close(errChan)
		return
	case <-exit_chan:
		log.Println("Exit by os.Interrupt syscall.SIGTERM")
		// initiate graceful shutdown
		close(closing)
		wg.Wait() //wait childs
		close(errChan)
		return
	}
}
