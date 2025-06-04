package config

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	// Project environment mode.
	// Available modes: "disable", "local", "dev", "prod".
	Env string `yaml:"env" env-required:"true"`

	Database `yaml:"database" env-required:"true"`

	HTTPServer `yaml:"http_server" env-required:"true"`

	TelegramBot `yaml:"telegram_bot" env-required:"true"`
}

type Database struct {
	// Path to the directory where migrations are stored.
	// Default path: ./migrations
	MigrationsPath string `yaml:"migrations_path" env-default:"./migrations"`

	StoragePath string `yaml:"storage_path" env-required:"true"`
}

type HTTPServer struct {
	Host        string        `yaml:"host" env-required:"true"`
	Port        int           `yaml:"port" env-required:"true"`
	Timeout     time.Duration `yaml:"timeout" env-required:"true"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-required:"true"`
}

type TelegramBot struct {
	Token     string `yaml:"token" env-required:"true"`
	ChannelID int64  `yaml:"channel_id"`

	// The size of one chunk stored in a telegram channel.
	// The default value is 15 Mb (not recommended to change)
	ChunkSize int64 `yaml:"chunk_size" env-default:"15728640"`
}

func New() (*Config, error) {
	return LoadConfigFromPath(GetConfigPath())
}

func GetConfigPath() string {
	var configPath string
	flag.StringVar(&configPath, "config", "", "Path to config file")
	flag.Parse()

	if configPath == "" {
		configPath = os.Getenv("CONFIG_PATH")
	}

	return configPath
}

func LoadConfigFromPath(configPath string) (*Config, error) {
	cfg := new(Config)

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, errors.New("failed to load config: config file not found")
	}

	if err := cleanenv.ReadConfig(configPath, cfg); err != nil {
		return nil, fmt.Errorf("fail to read config file: %s", err)
	}

	return cfg, nil
}
