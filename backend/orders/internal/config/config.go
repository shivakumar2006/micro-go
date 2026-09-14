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

type DBConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
}

type JWTConfig struct {
	AccessSecret  string `yaml:"access_secret"`
	RefreshSecret string `yaml:"refresh_secret"`
}

type CartConfig struct {
	URL string `yaml:"url"`
}

type Config struct {
	Env                string       `yaml:"env"`
	Server             ServerConfig `yaml:"server"`
	Cart               CartConfig   `yaml:"cart"`
	DB                 DBConfig     `yaml:"db"`
	JWT                JWTConfig    `yaml:"jwt"`
	InternalServiceKey string       `yaml:"internal_service_key"`
}

func LoadConfig() *Config {
	var configPath string

	configPath = os.Getenv("CONFIG_PATH")

	if configPath == "" {
		flags := flag.String("config", "", "path to config file")
		flag.Parse()
		configPath = *flags

		if configPath == "" {
			log.Fatal("config file path not found")
		}
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file doesn't exist : %v", err)
	}

	var config Config

	if err := cleanenv.ReadConfig(configPath, &config); err != nil {
		log.Fatalf("failed to load config file : %v", err)
	}

	if value := os.Getenv("DB_PASSWORD"); value != "" {
		config.DB.Password = value
	}

	if value := os.Getenv("JWT_ACCESS_SECRET"); value != "" {
		config.JWT.AccessSecret = value
	}

	if value := os.Getenv("JWT_REFRESH_SECRET"); value != "" {
		config.JWT.RefreshSecret = value
	}

	if value := os.Getenv("INTERNAL_SERVICE_KEY"); value != "" {
		config.InternalServiceKey = value
	}

	return &config
}
