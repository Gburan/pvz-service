package config

import (
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Server ServerConfig `yaml:"server"`
	DB     DBConfig     `yaml:"postgres"`
	App    AppConfig    `yaml:"app"`
}

type AppConfig struct {
	AllowedCategories []string `yaml:"allowed_categories" env:"ALLOWED_CATEGORIES" env-default:"электроника,одежда,обувь"`
	AllowedCities     []string `yaml:"allowed_cities" env:"ALLOWED_CITIES" env-default:"Москва,Санкт-Петербург,Казань"`
	AllowedUsers      []string `yaml:"allowed_users" env:"ALLOWED_USERS" env-default:"moderator,employee"`
	JWTToken          string   `yaml:"jwt_token" env:"JWT_TOKEN" env-default:""`
}

type ServerConfig struct {
	Address string `yaml:"address" env:"SERVER_ADDRESS" env-default:":8080"`
}

type DBConfig struct {
	MigrationsDir string `yaml:"migrations_dir" env:"MIGRATIONS_DIR" env-default:"./migrations"`
	Conn          string `yaml:"conn" env:"POSTGRES_CONN" env-default:""`
}

func MustLoad(configPath string) Config {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatal("Cannot find config file")
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatal("Error while reading config")
	}

	return cfg
}
