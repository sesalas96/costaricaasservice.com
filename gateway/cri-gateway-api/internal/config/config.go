// Package config carga config/{env}.yaml con overrides por env vars.
// El gateway no tiene DB; el config carga upstreams, claves y TTLs.
package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"

	"github.com/devsebas/saascr/gateway/cri-gateway-api/internal/upstream"
)

type Server struct {
	ReadTimeout     time.Duration `mapstructure:"read_timeout"`
	WriteTimeout    time.Duration `mapstructure:"write_timeout"`
	IdleTimeout     time.Duration `mapstructure:"idle_timeout"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"`
}

type Auth struct {
	PublicPEMPath  string   `mapstructure:"public_pem_path"`
	PublicPrefixes []string `mapstructure:"public_prefixes"`
}

type Revocation struct {
	IdentityURL  string        `mapstructure:"identity_url"`
	PollInterval time.Duration `mapstructure:"poll_interval"`
	Realms       []string      `mapstructure:"realms"`
}

type RateLimit struct {
	Rate    float64       `mapstructure:"rate"`     // req/sec por clave
	Burst   float64       `mapstructure:"burst"`    // capacidad
	IdleTTL time.Duration `mapstructure:"idle_ttl"` // limpieza
}

type Config struct {
	Env             string           `mapstructure:"env"`
	ServiceName     string           `mapstructure:"service_name"`
	Port            string           `mapstructure:"port"`
	LogLevel        string           `mapstructure:"log_level"`
	BaseHost        string           `mapstructure:"base_host"`
	UpstreamTimeout time.Duration    `mapstructure:"upstream_timeout"`
	Server          Server           `mapstructure:"server"`
	Auth            Auth             `mapstructure:"auth"`
	Revocation      Revocation       `mapstructure:"revocation"`
	RateLimit       RateLimit        `mapstructure:"rate_limit"`
	Upstreams       []upstream.Spec  `mapstructure:"upstreams"`
	ShutdownTimeout time.Duration    `mapstructure:"-"`
}

// Load lee el archivo config/{ENV}.yaml y aplica overrides por env vars.
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
	v.SetDefault("port", "18000")
	v.SetDefault("log_level", "info")
	v.SetDefault("upstream_timeout", "30s")
	v.SetDefault("server.read_timeout", "10s")
	v.SetDefault("server.write_timeout", "30s")
	v.SetDefault("server.idle_timeout", "60s")
	v.SetDefault("server.shutdown_timeout", "15s")
	v.SetDefault("revocation.poll_interval", "30s")
	v.SetDefault("rate_limit.rate", 50.0)   // 50 req/sec por clave
	v.SetDefault("rate_limit.burst", 100.0)
	v.SetDefault("rate_limit.idle_ttl", "5m")

	if err := v.ReadInConfig(); err != nil {
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
