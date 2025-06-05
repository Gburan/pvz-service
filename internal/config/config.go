package config

import (
	"log/slog"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Server ServerConfig `yaml:"server"`
	DB     DBConfig     `yaml:"postgres"`
	App    AppConfig    `yaml:"app"`
}

type AppConfig struct {
	Validation Validation
	Logging    Logging
	JWTToken   string `yaml:"jwt_token" env:"JWT_TOKEN" env-default:""`
}

type Logging struct {
	Output string `yaml:"output" env:"OUTPUT" env-default:"stdout"`
	Level  string `yaml:"level" env:"LEVEL" env-default:"debug"`
}

type Validation struct {
	AllowedCategories []string `yaml:"allowed_categories" env:"ALLOWED_CATEGORIES" env-default:"электроника,одежда,обувь"`
	AllowedCities     []string `yaml:"allowed_cities" env:"ALLOWED_CITIES" env-default:"Москва,Санкт-Петербург,Казань"`
	AllowedUsers      []string `yaml:"allowed_users" env:"ALLOWED_USERS" env-default:"moderator,employee"`
}

type ServerConfig struct {
	Rest RestConfig `yaml:"rest"`
	GRPC GRPCConfig `yaml:"grpc"`
}

type RestConfig struct {
	Address      string           `yaml:"address" env:"REST_ADDRESS" env-default:":8080"`
	Connsettings RestConnSettings `yaml:"conn_settings"`
}

type RestConnSettings struct {
	ReadTimeout  time.Duration `yaml:"read_timeout" env:"READ_TIMEOUT" env-default:"5s"`
	WriteTimeout time.Duration `yaml:"write_timeout" env:"WRITE_TIMEOUT" env-default:"5s"`
	IdleTimeout  time.Duration `yaml:"idle_timeout" env:"IDLE_TIMEOUT" env-default:"5m"`
}

type GRPCConfig struct {
	Address      string           `yaml:"address" env:"GRPC_ADDRESS" env-default:":3000"`
	ConnSettings GRPCConnSettings `yaml:"conn_settings"`
}

type GRPCConnSettings struct {
	MaxConnIdle time.Duration `yaml:"max_conn_idle" env:"GRPC_MAX_CONN_IDLE" env-default:"5m"`
	MaxConnAge  time.Duration `yaml:"max_conn_age" env:"GRPC_MAX_CONN_AGE" env-default:"10m"`
}

type DBConfig struct {
	MigrationsDir string       `yaml:"migrations_dir" env:"MIGRATIONS_DIR" env-default:"./migrations"`
	Conn          string       `yaml:"conn" env:"POSTGRES_CONN" env-default:""`
	PoolSettings  PoolSettings `yaml:"pool_settings"`
}

type PoolSettings struct {
	MaxConns        int32         `yaml:"max_conns" env:"MAX_CONNS" env-default:"15"`
	MinIdleConns    int32         `yaml:"min_idle_conns" env:"MIN_IDLE_CONNS" env-default:"5"`
	MaxConnIdleTime time.Duration `yaml:"max_conn_idle_time" env:"MAX_CONN_IDLE_TIME" env-default:"5m"`
	MaxConnLifetime time.Duration `yaml:"max_conn_lifetime" env:"MAX_CONN_LIFETIME" env-default:"10m"`
}

func MustLoad(configPath string) Config {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		slog.Error("Cannot find config file")
		os.Exit(1)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		slog.Error("Error while reading config")
		os.Exit(1)
	}

	return cfg
}
