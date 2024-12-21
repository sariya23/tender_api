package config

import (
	"flag"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type AppConfig struct {
	ServerAddress    string `env:"SERVER_ADDRESS"`
	PostgresConn     string `env:"POSTGRESS_CONN"`
	PosthresJDBC_URL string `env:"POSTGRES_JDBC_CONN"`
	PostgresUsername string `env:"POSTGRES_USERNAME"`
	PostgresPassword string `env:"POSTRGRES_PASSWORD"`
	PostgresHost     string `env:"POSTGRES_HOST"`
	PostgresPort     int    `env:"POSTGRES_PORT"`
	PostgresDatabase string `env:"POSTGRES_DB"`
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
