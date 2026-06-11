package config

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type ServerConfig struct {
	Addr string `yaml:"address" env-required:"true"`
}

type DBConfig struct {
	Host     string `yaml:"host" env-required:"true"`
	Port     int    `yaml:"port" env-required:"true"`
	User     string `yaml:"user" env-required:"true"`
	Password string `yaml:"password" env-required:"true"`
	DBName   string `yaml:"dbname" env-required:"true"`
}

type JWTConfig struct {
	AccessSecret  string        `yaml:"access_secret" env-required:"true"`
	RefreshSecret string        `yaml:"refresh_secret" env-required:"true"`
	AccessExpiry  time.Duration `yaml:"access_expiry" env-required:"true"`
	RefreshExpiry time.Duration `yaml:"refresh_expiry" env-required:"true"`
}

type Config struct {
	Env    string       `yaml:"env" env-required:"true"`
	Server ServerConfig `yaml:"server"`
	DB     DBConfig     `yaml:"db"`
	JWT    JWTConfig    `yaml:"jwt"`
}

func LoadConfig() *Config {
	var configPath string

	configPath = os.Getenv("CONFIG_PATH")

	if configPath == "" {
		flags := flag.String("config", "", "path to config file") // it takes name, value and usage
		flag.Parse()

		configPath = *flags

		if configPath == "" {
			log.Fatal("config path is not set")
		}
	}

	// os.Stat returns the file info if info is not present means there's an error
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file doesn't exist %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("Failed to read config file %s", err.Error())
	}

	log.Println("Config added successfully")

	return &cfg
}
