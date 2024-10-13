package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type HttpServer struct {
	Addr string `yaml:"address" env-required:"true`
}
type Config struct {
	Env          string `yaml:"env" env:"Env" env_required:"true" env-default:"production"`
	Storage_Path string `yaml:"storage_path" env_required:"true"`
	HttpServer   `yaml:"http_server"`
}

// create an function to used above struct

func MustLoad() *Config {
	var config_path string

	config_path = os.Getenv("CONFIG_PATH")

	if config_path == "" {
		flags := flag.String("config", "", "path to the configure file")
		flag.Parse()
		config_path = *flags

		if config_path == "" {
			log.Fatal("Config Path is not set")
		}
	}

	if _, err := os.Stat(config_path); os.IsNotExist(err) {
		log.Fatalf("Config file does not exist %s", config_path)
	}
	var cfg Config
	err := cleanenv.ReadConfig(config_path, &cfg)

	if err != nil {
		log.Fatalf("Error loading config file: %s", err.Error())
	}
	return &cfg
}
