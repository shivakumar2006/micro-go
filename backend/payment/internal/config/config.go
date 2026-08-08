package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type ServerConfig struct {
	Addr string `yaml:"address"`
}

type DBConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
}

type JWTConfig struct {
	AccessTokenSecret  string `yaml:"access_secret"`
	RefreshTokenSecret string `yaml:"refresh_secret"`
}

type OrdersConfig struct {
	URL string `yaml:"url"`
}

type Config struct {
	Env    string       `yaml:"env"`
	Server ServerConfig `yaml:"server"`
	Orders OrdersConfig `yaml:"orders"`
	DB     DBConfig     `yaml:"db"`
	JWT    JWTConfig    `yaml:"jwt"`
}

func LoadConfig() *Config {
	// load env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf(".env file is not found : %v", err)
	}

	var configPath string

	configPath = os.Getenv("CONFIG_PATH")

	if configPath == "" {
		flags := flag.String("config", "", "path to config file")
		flag.Parse()
		configPath = *flags

		if configPath == "" {
			log.Fatal("Config file path is not set")
		}
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatal("Config file does not exist: ", configPath)
	}

	var config Config

	if err := cleanenv.ReadConfig(configPath, &config); err != nil {
		log.Fatal("Error reading config: ", err)
	}

	log.Println("config successfully added")

	return &config
}
