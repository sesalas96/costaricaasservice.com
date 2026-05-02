// Package config carga config/{env}.yaml con overrides por env vars.
// Convención: ENV var selecciona el archivo, las demás overridean campos.
package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Database struct {
	URL      string `mapstructure:"url"`
	MaxConns int    `mapstructure:"max_conns"`
	MinConns int    `mapstructure:"min_conns"`
}

type Server struct {
	ReadTimeout     time.Duration `mapstructure:"read_timeout"`
	WriteTimeout    time.Duration `mapstructure:"write_timeout"`
	IdleTimeout     time.Duration `mapstructure:"idle_timeout"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"`
}

type Upstreams struct {
	HaciendaURL string `mapstructure:"hacienda_url"`
	CCSSURL     string `mapstructure:"ccss_url"`
	AuditURL    string `mapstructure:"audit_url"`
}

type CORS struct {
	AllowedOrigins []string `mapstructure:"allowed_origins"`
}

type Config struct {
	Env             string        `mapstructure:"env"`
	ServiceName     string        `mapstructure:"service_name"`
	Port            string        `mapstructure:"port"`
	LogLevel        string        `mapstructure:"log_level"`
	BaseHost        string        `mapstructure:"base_host"`
	Database        Database      `mapstructure:"database"`
	Server          Server        `mapstructure:"server"`
	Upstreams       Upstreams     `mapstructure:"upstreams"`
	CORS            CORS          `mapstructure:"cors"`
	ShutdownTimeout time.Duration `mapstructure:"-"`
}

// Load lee el archivo config/{ENV}.yaml y aplica overrides por env vars.
// El archivo se busca relativo al cwd y luego en ./config/.
func Load() (*Config, error) {
	env := os.Getenv("ENV")
	if env == "" {
		env = "local"
	}

	v := viper.New()
	v.SetConfigName(env)
	v.SetConfigType("yaml")
	v.AddConfigPath("./config")
	v.AddConfigPath(".")

	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// defaults
	v.SetDefault("env", env)
	v.SetDefault("port", "8080")
	v.SetDefault("log_level", "info")
	v.SetDefault("server.read_timeout", "10s")
	v.SetDefault("server.write_timeout", "30s")
	v.SetDefault("server.idle_timeout", "60s")
	v.SetDefault("server.shutdown_timeout", "15s")
	v.SetDefault("database.max_conns", 10)
	v.SetDefault("database.min_conns", 2)

	if err := v.ReadInConfig(); err != nil {
		// archivo opcional — config solo de env vars también es válido
		var notFound viper.ConfigFileNotFoundError
		if !asErr(err, &notFound) {
			return nil, fmt.Errorf("read config: %w", err)
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}
	cfg.ShutdownTimeout = cfg.Server.ShutdownTimeout
	return &cfg, nil
}

func asErr(err error, target any) bool {
	type asInterface interface {
		As(any) bool
	}
	if x, ok := err.(asInterface); ok {
		return x.As(target)
	}
	return false
}
