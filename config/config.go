package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	//"gitlab.boquar.tech/galileosky/pkg/acl"
)

type (
	Config struct {
		App   `yaml:"app"`
		Log   `yaml:"log"`
		PG    `yaml:"pg"`
		Trace `yaml:"trace"`
		HTTP  `yaml:"http"`
		//Acl      *acl.ConfigAcl `yaml:"acl"`
		GRPC     `yaml:"grpc"`
		KeyCloak `yaml:"keycloak"`
	}

	App struct {
		Name        string `yaml:"name" env:"APP_NAME"`
		Environment string `yaml:"environment" env:"APP_ENVIRONMENT" env-default:"develop"`
	}

	Log struct {
		//Level string `yaml:"level" env:"LOG_LEVEL" env-default:"info"`
		Level string `env-required:"true" yaml:"level"   env:"LOG_LEVEL"`
	}

	PG struct {
		RSID        string `yaml:"rsid" env:"RSID"`
		PoolMax     int    `yaml:"pool_max" env:"PG_POOL_MAX"`
		DatabaseURL string `yaml:"database_url" env:"PG_DATABASE_URL"`
		SourceURL   string `yaml:"source_url" env:"PG_DATABASE_URL"`
	}

	Trace struct {
		URL   string `yaml:"url" env:"TRACER_URL" env-default:"http://localhost:14268/api/traces?service=customer-administration"`
		Level int    `env-required:"true" yaml:"level"   env:"TRACE_LEVEL"`
	}

	HTTP struct {
		Port string `yaml:"port" env:"HTTP_PORT"`
	}

	GRPC struct {
		Port string `yaml:"port" env:"GRPC_PORT"`
	}

	KeyCloak struct {
		URL          string `yaml:"url" env:"KEYCLOAK_URL" env-default:"https://tpm-keycloak.boquar.tech"`
		UserName     string `yaml:"username" env:"KEYCLOAK_USERNAME" env-default:"acladminservice"`
		ClientId     string `yaml:"clientid" env:"KEYCLOAK_CLIENTID" env-default:"aclaAdmin"`
		ClientSecret string `yaml:"clientsecret" env:"KEYCLOAK_CLIENTSECRET" env-default:"Freedom123_!"`
		Timeout      int    `yaml:"timeout" env:"KEYCLOAK_TIMEOUT"`
		Realm        string `yaml:"realm" env:"KEYCLOAK_REALM"`
	}
)

func NewConfig() (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig("config/config.yml", cfg)
	if err != nil {
		return nil, err
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
