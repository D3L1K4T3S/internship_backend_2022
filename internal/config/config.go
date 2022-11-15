package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"internship_bachend_2022/pkg/logging"
	"sync"
)

type Config struct {
	IsDebug *bool `yaml:"is_debug" env-default:"true"`
	Listen  struct {
		Type   string `yaml:"type" env-default:"port"`
		BindIP string `yaml:"bind_ip" env-default:"127.0.0.1"`
		Port   string `yaml:"port" env-default:"8080"`
	} `yaml:"listen"`
	Storage StorageConfig `yaml:"storage"`
}

type StorageConfig struct {
	Host     string `json:"host,omitempty"`
	Port     string `json:"port,omitempty"`
	Database string `json:"database,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

var instances *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		logger := logging.GetLogger()
		logger.Infof("read applicatiob configuration")
		instances = &Config{}
		if err := cleanenv.ReadConfig("config.yml", instances); err != nil {
			help, _ := cleanenv.GetDescription(instances, nil)
			logger.Info(help)
			logger.Fatal(err)
		}
	})
	return instances
}
