package config

import (
	"flag"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env    string         `yaml:"env" env-default:"local"`
	Server ServerConfig   `yaml:"server"`
	DB     DatabaseConfig `yaml:"database"`
}

type ServerConfig struct {
	Host string `yaml:"host" env-default:"localhost"`
	Port string `yaml:"port" env-default:"3002"`
}

type DatabaseConfig struct {
	Driver string `yaml:"driver"`
	Host   string `yaml:"host" env-default:"localhost"`
	Port   string `yaml:"port" env-default:"5432"`
	User   string `yaml:"user" env-default:"postgres"`
	Pass   string `yaml:"pass" env-default:"postgres"`
	Name   string `yaml:"name" env-default:"antisocial"`
}

func MustLoad() *Config {
	configPath := fetchConfigPath()
	if configPath == "" {
		panic("config path is empty")
	}
	return MustLoadPath(configPath)
}

func MustLoadPath(configPath string) *Config {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file does not exist: " + configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("cannot read config: " + err.Error())
	}

	return &cfg
}

func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
