package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"img-build-ci-runner/configs"

	"github.com/kirsle/configdir"
)

const appname = "img-build-ci-runner"
const cfgExample = "./configs/cfg_example.json"

type Config struct {
	path     string
	settings *AppSettings
}

func New() *Config {
	configPath := configdir.LocalConfig(appname)
	err := configdir.MakePath(configPath) // Ensure it exists.
	if err != nil {
		panic(err)
	}

	configFile := filepath.Join(configPath, "config.json")
	log.Printf("Config path: %s\n", configFile)
	var settings AppSettings

	if _, err = os.Stat(configFile); os.IsNotExist(err) {
		// Create the new config file
		// Write exaple into it
		fh, err := os.Create(configFile)
		if err != nil {
			log.Fatal(fmt.Sprintf("Can't create default config file. Path %s. Error: %v\n", configFile, err))
			panic(err)
		}
		defer fh.Close()

		if _, err = fh.WriteString(configs.Cfg_example); err != nil {
			log.Fatal(fmt.Sprintf("Can't write default configuration to config file. Path %s; Conf-on: %s. Error: %v\n", configFile, configs.Cfg_example, err))
			panic(err)
		}

		json.Unmarshal([]byte(configs.Cfg_example), &settings)
	} else {
		// Load the existing file.
		fh, err := os.Open(configFile)
		if err != nil {
			log.Fatal(fmt.Sprintf("Can't open config file. Path %s. Error: %v\n", configFile, err))
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

func (cfg *Config) UpdateSettings() error {
	fh, err := os.Open(cfg.path)
	if err != nil {
		return fmt.Errorf("Can't read config file %s and get settings. Error: %w\n", cfg.path, err)
	}
	defer fh.Close()

	decoder := json.NewDecoder(fh)
	decoder.Decode(&cfg.settings)
	return nil
}

func (cfg *Config) GetSettings(name string) string {
	switch name {
	case "AltSiteUrl":
		return cfg.settings.AltSiteUrl
	case "GiteaWfUrl":
		return cfg.settings.GiteaWfUrl
	case "GiteaRepoUrl":
		return cfg.settings.GiteaRepoUrl
	case "DepsCronGroupExp":
		return cfg.settings.DepsCronGroupExp
	case "VersCronGroupExp":
		return cfg.settings.VersCronGroupExp
	case "VersCheckImgGroup":
		return cfg.settings.VersCheckImgGroup
	case "DepsCronImgGroup":
		return cfg.settings.DepsCronImgGroup
	case "GiteaToken":
		return cfg.settings.GiteaToken
	default:
		err := fmt.Errorf("Setting %s is not in settings list. Check your code and settings list in config %s\n", name, cfg.path)
		log.Fatal(err)
		return ""
	}
}

func (cfg *Config) GetBranches() []string {
	if len(cfg.settings.Branches) <= 0 {
		err := fmt.Errorf("Setting Branches from config is empty. Set value to %s\n", cfg.path)
		log.Fatal(err)
		return nil
	}

	return cfg.settings.Branches
}
