package config

import "github.com/spf13/viper"

type Config struct {
	Port      string              `mapstructure:"HTTP_PORT"`
	SecretKey string              `mapstructure:"SECRET_KEY"`
	NoSQLdb   NoSQLDatabaseConfig `mapstructure:",squash"`
}

type NoSQLDatabaseConfig struct {
	Host string `mapstructure:"NOSQL_DATABASE_HOST"`
	Port string `mapstructure:"NOSQL_DATABASE_PORT"`
	Name string `mapstructure:"NOSQL_DATABASE_NAME"`
}

func Load() (Config, error) {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		return Config{}, nil
	}

	viper.BindEnv(
		"HTTP_PORT",
		"SECRET_KEY",
		"NOSQL_DATABASE_HOST",
		"NOSQL_DATABASE_PORT",
		"NOSQL_DATABASE_NAME",
	)

	var c Config
	err = viper.Unmarshal(&c)

	return c, err
}
