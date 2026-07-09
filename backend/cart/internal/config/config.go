package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type ServerConfig struct {
	Address string `yaml:"address"`
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

type Config struct {
	Env    string `yaml:"env"`
	Server ServerConfig
	DB     DBConfig
	JWT    JWTConfig
}

func LoadConfig() (*Config, error) {
	var ConfigPath string

	ConfigPath = os.Getenv("CONFIG_PATH")

	if ConfigPath == "" {
		flags := flag.String("config", "", "path to config file")
		flag.Parse()
		ConfigPath = *flags

		if ConfigPath == "" {
			log.Fatal("config file path not provided")
		}
	}

	if _, err := os.Stat(ConfigPath); os.IsNotExist(err) {
		log.Fatalf("Config file doesn't exist : %v", err)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(ConfigPath, &cfg); err != nil {
		panic("failed to load config")
	}

	log.Println("config successfully added")

	return &cfg, nil
}
