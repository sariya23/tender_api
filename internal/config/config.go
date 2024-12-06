package config

import (
	"flag"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type AppConfig struct {
	ServerAddress    string `yaml:"server_address"`
	PostgresConn     string `yaml:"postgress_conn"`
	PosthresJDBC_URL string `yaml:"postgres_jdbc_conn"`
	PostgresUsername string `yaml:"postgres_username"`
	PostgresPassword string `yaml:"postgres_password"`
	PostgresHost     string `yaml:"postgres_host"`
	PostgresPort     int    `yaml:"postgres_port"`
	PostgresDatabase string `yaml:"postgres_database"`
}

func MustLoad() *AppConfig {
	path := fetchConfigPath()
	if path == "" {
		panic("config path is empty")
	}
	return MustLoadByPath(path)
}

// Priority: flag > env > default
func fetchConfigPath() string {
	var path string
	flag.StringVar(&path, "config", "", "path to config file")
	flag.Parse()

	if path == "" {
		path = os.Getenv("CONFIG_PATH")
	}
	return path
}

func MustLoadByPath(path string) *AppConfig {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file does not exists: " + path)
	}

	var cfg AppConfig
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("cannot read config: " + err.Error())
	}
	return &cfg
}
