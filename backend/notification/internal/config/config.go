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

type KafkaConfig struct {
	Addr    string `yaml:"address"`
	Topic   string `yaml:"topic"`
	GroupID string `yaml:"group_id"`
}
type BrevoConfig struct {
	SMTPHost     string `yaml:"smtp_host"`
	SMTPPort     int    `yaml:"smtp_port"`
	SMTPUser     string `yaml:"smtp_user"`
	SMTPPassword string `yaml:"smtp_password"`
	SenderEmail  string `yaml:"sender_email"`
	SenderName   string `yaml:"sender_name"`
}

type Config struct {
	Server ServerConfig `yaml:"server"`
	Kafka  KafkaConfig  `yaml:"kafka"`
	Brevo  BrevoConfig  `yaml:"brevo"`
}

func LoadConfig() *Config {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("failed to load env : %v", err)
	}

	var configPath string

	configPath = os.Getenv("CONFIG_PATH")

	if configPath == "" {
		flags := flag.String("config", "", "path to config file")
		flag.Parse()
		configPath = *flags

		if configPath == "" {
			log.Fatalf("config path is not set")
		}
	}

	if _, err := os.Stat(configPath); err != nil {
		log.Fatalf("config file is not found: %v", err)
	}

	var config Config

	if err := cleanenv.ReadConfig(configPath, &config); err != nil {
		log.Fatalf("failed to read config: %v", err)
	}

	if err := cleanenv.ReadEnv(&config); err != nil {
		log.Fatalf("failed to read env: %v", err)
	}

	log.Println("config file loaded successfully")

	return &config
}
