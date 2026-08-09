package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type ServerConfig struct {
	Addr string `yaml:"address"`
}

type ServicesConfig struct {
	Auth    Service `yaml:"auth"`
	Vehicle Service `yaml:"vehicle"`
	Cart    Service `yaml:"cart"`
	Orders  Service `yaml:"orders"`
	Payment Service `yaml:"payment"`
}

type Service struct {
	URL string `yaml:"url"`
}

type JWTConfig struct {
	AccessSecret  string `yaml:"access-secret"`
	RefreshSecret string `yaml:"refresh-secret"`
}

type Config struct {
	Env      string         `yaml:"env"`
	Server   ServerConfig   `yaml:"server"`
	Services ServicesConfig `yaml:"services"`
	JWT      JWTConfig      `yaml:"jwt"`
}

func LoadConfig() *Config {
	var configPath string

	configPath = os.Getenv("CONFIG_PATH")

	if configPath == "" {
		flags := flag.String("config", " ", "path to config file")
		flag.Parse()
		configPath = *flags

		if configPath == "" {
			log.Fatal("config file path is not provided")
		}
	}

	if _, err := os.Stat(configPath); err != nil {
		log.Fatalf("Config file doesn't exist : %v", err)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("config file not found : %v", err)
	}

	log.Println("config loaded successfully")

	return &cfg
}
