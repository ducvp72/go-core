package config

import (
	"fmt"
	"log"
	"os"

	"path/filepath"

	"gopkg.in/ini.v1"
)

type ServiceConfig struct {
	Env  string
	App  string
	Name string
	Port int
}

var ServConfig ServiceConfig

func LoadConfig() {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "dev" // default to dev
	}

	workingDirpath, _ := os.Getwd()
	configPath := filepath.Join(workingDirpath, fmt.Sprintf("config/config.ini.%s", env))

	cfg, err := ini.Load(configPath)
	if err != nil {
		log.Fatalf("Failed to read config file %s: %v", configPath, err)
	}

	err = cfg.Section("Service").MapTo(&ServConfig)
	if err != nil {
		log.Fatalf("Failed to map config section: %v", err)
	}
}
