package configs

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type AppConfig struct {
	DotSSHPath    string   `yaml:"dot-ssh-path"`
	Notifications []string `yaml:"notifications"`
	SyncFilePaths map[string][]struct {
		ConfigName string `yaml:"config-name"`
		Dst        string `yaml:"dst"`
	} `yaml:"sync-file-paths"`
}

var cfg AppConfig

func init() {
	err := cleanenv.ReadConfig("config.yaml", &cfg)
	if err != nil {
		log.Println("failed to read config file config.yaml")
		panic(err)
	}
}

func GetAppConfig() AppConfig {
	return cfg
}
