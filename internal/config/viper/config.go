package config

import (
	"bytes"
	"fmt"
	"log"

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

	viper.AddConfigPath("/etc/appname/")
	viper.AddConfigPath(configPath)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			viper.ReadConfig(bytes.NewBuffer([]byte(configs.Cfg_example)))
			viper.WriteConfig() // writes current config to predefined path set by 'viper.AddConfigPath()' and 'viper.SetConfigName'
			viper.SafeWriteConfig()
		} else {
			log.Fatalln(err)
		}
	}

	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
	})
	v.WatchConfig()

	return &Config{
		v,
	}
}
