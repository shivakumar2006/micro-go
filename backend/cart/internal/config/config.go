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

type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
	Db       int    `yaml:"db"`
}

type JWTConfig struct {
	AccessSecret  string `yaml:"access_secret"`
	RefreshSecret string `yaml:"refresh_secret"`
}

type VehicleConfig struct {
	URL string `yaml:"url"`
}

type Config struct {
	Env     string        `yaml:"env"`
	Server  ServerConfig  `yaml:"server"`
	Vechile VehicleConfig `yaml:"vehicle"`
	Redis   RedisConfig   `yaml:"redis"`
	DB      DBConfig      `yaml:"db"`
	JWT     JWTConfig     `yaml:"jwt"`
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

	if value := os.Getenv("DB_PASSWORD"); value != "" {
		cfg.DB.Password = value
	}

	if value := os.Getenv("JWT_ACCESS_SECRET"); value != "" {
		cfg.JWT.AccessSecret = value
	}

	if value := os.Getenv("JWT_REFRESH_SECRET"); value != "" {
		cfg.JWT.RefreshSecret = value
	}

	if value := os.Getenv("REDIS_PASSWORD"); value != "" {
		cfg.Redis.Password = value
	}

	log.Println("config successfully added")

	return &cfg, nil
}
