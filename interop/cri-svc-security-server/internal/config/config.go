// Package config carga config/{env}.yaml para el cri-svc-security-server.
package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Server struct {
	ReadTimeout     time.Duration `mapstructure:"read_timeout"`
	WriteTimeout    time.Duration `mapstructure:"write_timeout"`
	IdleTimeout     time.Duration `mapstructure:"idle_timeout"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"`
}

type Member struct {
	Realm             string `mapstructure:"realm"`
	Slug              string `mapstructure:"slug"`
	PrivatePEMPath    string `mapstructure:"private_pem_path"`
	MemberDispatchURL string `mapstructure:"member_dispatch_url"`
}

type Hub struct {
	URL      string        `mapstructure:"url"`
	CacheTTL time.Duration `mapstructure:"cache_ttl"`
}

type Audit struct {
	URL string `mapstructure:"url"` // ej: http://localhost:18102 — vacío = audit deshabilitado
}

type Peer struct {
	Member        string `mapstructure:"member"`
	InboxURL      string `mapstructure:"inbox_url"`
	PublicKeyPath string `mapstructure:"public_key_path"` // opcional: anchor local de pubkey
}

// PeerKeyPath devuelve la ruta de la pubkey local de un peer (si está configurada).
func (c *Config) PeerKeyPath(member string) string {
	for _, p := range c.Peers {
		if p.Member == member {
			return p.PublicKeyPath
		}
	}
	return ""
}

type Config struct {
	Env             string        `mapstructure:"env"`
	ServiceName     string        `mapstructure:"service_name"`
	Port            string        `mapstructure:"port"`
	LogLevel        string        `mapstructure:"log_level"`
	BaseHost        string        `mapstructure:"base_host"`
	CallTimeout     time.Duration `mapstructure:"call_timeout"`
	Member          Member        `mapstructure:"member"`
	Hub             Hub           `mapstructure:"hub"`
	Audit           Audit         `mapstructure:"audit"`
	Peers           []Peer        `mapstructure:"peers"`
	Server          Server        `mapstructure:"server"`
	ShutdownTimeout time.Duration `mapstructure:"-"`
}

// PeerByMember devuelve el inbox URL configurado de un member.
func (c *Config) PeerByMember(member string) (string, bool) {
	for _, p := range c.Peers {
		if p.Member == member {
			return p.InboxURL, true
		}
	}
	return "", false
}

func Load() (*Config, error) {
	env := os.Getenv("ENV")
	if env == "" {
		env = "local"
	}
	// CONFIG override permite correr varias instancias del SS (una por member)
	// con archivos separados: config/local-hacienda.yaml, config/local-registro-civil.yaml.
	configName := os.Getenv("CONFIG")
	if configName == "" {
		configName = env
	}

	v := viper.New()
	v.SetConfigName(configName)
	v.SetConfigType("yaml")
	v.AddConfigPath("./config")
	v.AddConfigPath(".")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	v.SetDefault("env", env)
	v.SetDefault("port", "18200")
	v.SetDefault("log_level", "info")
	v.SetDefault("call_timeout", "5s")
	v.SetDefault("hub.cache_ttl", "5m")
	v.SetDefault("server.read_timeout", "10s")
	v.SetDefault("server.write_timeout", "30s")
	v.SetDefault("server.idle_timeout", "60s")
	v.SetDefault("server.shutdown_timeout", "15s")

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
