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

type Config struct {
	Env    string       `yaml:"env"`
	Server ServerConfig `yaml:"server"`
	DB     DBConfig     `yaml:"db"`
	Redis  RedisConfig  `yaml:"redis"`
	JWT    JWTConfig    `yaml:"jwt"`
}

func LoadConfig() (*Config, error) {
	var ConfigPath string

	ConfigPath = os.Getenv("CONFIG_PATH")

	if ConfigPath == "" {
		flags := flag.String("config", "", "path to config file")
		flag.Parse()
		ConfigPath = *flags

		if ConfigPath == "" {
			log.Fatal("config path is not set")
		}
	}

	if _, err := os.Stat(ConfigPath); os.IsNotExist(err) {
		log.Fatalf("config file doesn't exist %s", ConfigPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(ConfigPath, &cfg); err != nil {
		panic("failed to load config")
	}

	log.Println("config added successfully")

	return &cfg, nil
}
