package config

import (
	"bytes"
	"log"
	"path/filepath"

	"img-build-ci-runner/configs"

	"github.com/fsnotify/fsnotify"
	"github.com/kirsle/configdir"
	"github.com/spf13/viper"
)

const appname = "img-build-ci-runner"

type Config struct {
	*viper.Viper
}

func New() *Config {
	v := viper.New()
	v.SetConfigName("config") // <- имя конфигурационного файла
	v.SetConfigType("json")

	configPath := configdir.LocalConfig(appname)
	err := configdir.MakePath(configPath) // Ensure it exists.
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("Config path: %s\n", configPath)

	v.AddConfigPath(configPath)
	v.AddConfigPath(filepath.Join("/etc/", appname))

	if err := v.ReadInConfig(); err != nil {
		log.Println("Config file is not found. Error: ", err)
		log.Println("Create config from template...")
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			v.ReadConfig(bytes.NewBuffer([]byte(configs.Cfg_example)))
			err = v.WriteConfig() // writes current config to predefined path set by 'viper.AddConfigPath()' and 'viper.SetConfigName'
			if err != nil {
				log.Println("Can't save config to predefined file path")
			}
			v.SafeWriteConfig()
		} else {
			log.Fatalln(err)
		}
	}

	v.OnConfigChange(func(e fsnotify.Event) {
		log.Println("Config file changed:", e.Name)
	})
	v.WatchConfig()

	return &Config{
		v,
	}
}
