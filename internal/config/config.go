package config

import (
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/lunn06/video-hosting/internal/models"
)

var CFG Config

type Config struct {
	HTTPServer `yaml:"http_server"`
	Database   `yaml:"database"`
}

type HTTPServer struct {
	Address string `yaml:"address" env-default:"127.0.0.1"`
	Port    string `yaml:"port" env-default:"8080"`
}

type Database struct {
	Address  string `yaml:"address"`
	User     string `yaml:"user"`
	Name     string `yaml:"name"`
	Password string `yaml:"password"`
	SSLMode  string `yaml:"ssl_mode"`
}

type DatabaseDefaults struct {
	Roles []models.Role
}

func Init() {
	CFG = MustLoad("configs/main.yaml")
}

func MustLoad(configPath string) Config {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}
	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return cfg
}

func MustLoadDatabaseDefaults(configPath string) DatabaseDefaults {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var defaults DatabaseDefaults

	if err := cleanenv.ReadConfig(configPath, &defaults); err != nil {
		log.Fatalf("MustLoadDatabaseDefaults() error = %v", err)
	}

	return defaults
}
