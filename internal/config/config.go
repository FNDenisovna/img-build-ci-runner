package config

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	model "altpack-vers-checker/internal/integration/model"

	"github.com/kirsle/configdir"
)

const appname = "altpack-vers-checker"
const imgPkgFilePath = "/var/lib/altpack-vers-checker/image_list.json"

type Config struct {
	path     string
	settings *model.AppSettings
}

func New() *Config {
	configPath := configdir.LocalConfig(appname)
	err := configdir.MakePath(configPath) // Ensure it exists.
	if err != nil {
		panic(err)
	}

	configFile := filepath.Join(configPath, "config.json")

	var settings model.AppSettings

	if _, err = os.Stat(configFile); os.IsNotExist(err) {
		// Create the new config file.
		fh, err := os.Create(configFile)
		if err != nil {
			panic(err)
		}
		defer fh.Close()

		encoder := json.NewEncoder(fh)
		encoder.Encode(&settings)
	} else {
		// Load the existing file.
		fh, err := os.Open(configFile)
		if err != nil {
			panic(err)
		}
		defer fh.Close()

		decoder := json.NewDecoder(fh)
		decoder.Decode(&settings)
	}
	return &Config{
		path:     configFile,
		settings: &settings,
	}
}

func (cfg *Config) UpdateCfgSettings() error {
	fh, err := os.Open(cfg.path)
	if err != nil {
		return nil, fmt.Errorf("Can't read config file %s and get settings. Error: %w", cfg.path, err)
	}
	defer fh.Close()

	decoder := json.NewDecoder(fh)
	decoder.Decode(&cfg.settings)
	return nil
}

func (cfg *Config) GetBranches() []string {
	return cfg.settings.Branches
}

func (cfg *Config) GetImgPkgList() ([]model.ImgPkg, error) {
	if err := downloadFile(cfg.settings.ImgPkgFileUrl, imgPkgFilePath); err != nil {
		return nil, fmt.Errorf("Can't get file from %s with image-package list. Error: %w", cfg.settings.ImgPkgFileUrl, err)
	}

	imgPkgFile, err := os.Open(imgPkgFilePath)
	if err != nil {
		return nil, fmt.Errorf("Can't open file with image-package list for checking version: %s. Error: %w", imgPkgFilePath, err)
	}
	defer imgPkgFile.Close()

	var res []model.ImgPkg
	decoder := json.NewDecoder(imgPkgFile)
	decoder.Decode(&res)
	return res, nil
}

func downloadFile(url, destPath string) error {
	// Create the file
	out, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
