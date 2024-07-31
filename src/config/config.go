package config

import (
	"crypto/sha256"
	"fmt"
	"os"

	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v3"
)

type Config struct {
	MaxConnections int    `yaml:"maxConnections" envconfig:"MAX_CONNECTIONS" validate:"min=1"`
	MaxCacheSize   int    `yaml:"maxCacheSize" envconfig:"MAX_CACHE_SIZE" validate:"min=1024"`
	BufferSize     int    `yaml:"bufferSize" envconfig:"BUFFER_SIZE" validate:"min=1024"`
	DataDir        string `yaml:"dataDir" envconfig:"DATA_DIR"`
	Port           string `yaml:"port" envconfig:"SERVER_PORT"`
	Passphrase     string `yaml:"passphrase" envconfig:"PASSPHRASE"`
	User           string `yaml:"user" envconfig:"USER"`
	NoAuth         bool   `yaml:"noAuth" envconfig:"NO_AUTH"`
	Version        string `yaml:"version" envconfig:"VERSION"`
	ConfigDir      string `yaml:"configDir" envconfig:"CONFIG_DIR"`
}

func LoadConfig() (Config, error) {
	cfg := Config{}
	yamlConfig(&cfg)
	envConfig(&cfg)
	if cfg.MaxConnections == 0 {
		cfg.MaxConnections = 10
	}
	if cfg.MaxCacheSize == 0 {
		cfg.MaxCacheSize = 2048
	}
	if cfg.BufferSize == 0 {
		cfg.BufferSize = 1024
	}
	if cfg.DataDir == "" {
		cfg.DataDir = "./data"
	}
	if cfg.ConfigDir == "" {
		cfg.ConfigDir = "./config"
	}
	if cfg.Port == "" {
		cfg.Port = "8080"
	}
	if cfg.Passphrase == "" {
		h := sha256.New()
		h.Write([]byte("password"))
		bs := h.Sum(nil)
		//Hash the passphrase to store it securely in the config
		cfg.Passphrase = string(bs)
	}
	if cfg.User == "" {
		cfg.User = "mew"
	}
	if cfg.Version == "" {
		cfg.Version = "1.0.0"
	}
	return cfg, nil
}

func processError(err error) {
	fmt.Println(err)
	os.Exit(2)
}

func yamlConfig(cfg *Config) {
	if _, err := os.Stat("config.yml"); os.IsNotExist(err) {
		return
	}
	f, err := os.Open("config.yml")
	if err != nil {
		processError(err)
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(cfg)
	if err != nil {
		processError(err)
	}
}

func envConfig(cfg *Config) {
	err := envconfig.Process("", cfg)
	if err != nil {
		processError(err)
	}
}
