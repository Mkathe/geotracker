package config

import (
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

type Config struct {
	ConnStr string `mapstructure:"CONN_DB_POSTGRES"`
	Port    string `mapstructure:"PORT"`
	Ip      string `mapstructure:"IP"`
	//ConnStrOutside string `mapstructure:"CONN_DB_POSTGRES_OUTSIDE"`
	//MigrationsPath string `mapstructure:"MIGRATION_PATH"`
	//Brokers        string `mapstructure:"KAFKA_BROKERS"`
	//Topic          string `mapstructure:"KAFKA_TOPIC"`
	//KeycloakURL    string `mapstructure:"KEYCLOAK_URL"`
	//KeycloakRole   string `mapstructure:"KEYCLOAK_ROLE"`
}

var cfg *Config

func Load() error {
	v := viper.New()
	v.AutomaticEnv()

	if err := v.BindEnv("CONN_DB_POSTGRES"); err != nil {
		return err
	}

	if err := v.BindEnv("IP"); err != nil {
		return err
	}

	if err := v.BindEnv("PORT"); err != nil {
		return err
	}

	if err := v.BindEnv("CONN_DB_POSTGRES_OUTSIDE"); err != nil {
		return err
	}

	if err := v.BindEnv("MIGRATION_PATH"); err != nil {
		return err
	}

	if err := v.BindEnv("KAFKA_BROKERS"); err != nil {
		return err
	}

	if err := v.BindEnv("KAFKA_TOPIC"); err != nil {
		return err
	}

	if err := v.BindEnv("KEYCLOAK_URL"); err != nil {
		return err
	}

	if err := v.BindEnv("KEYCLOAK_ROLE"); err != nil {
		return err
	}

	if err := v.Unmarshal(&cfg); err != nil {
		return err
	}

	return nil
}

func Get() *Config {
	return cfg
}
