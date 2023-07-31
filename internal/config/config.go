package config

import "github.com/spf13/viper"

type Config struct {
	Port               string             `mapstructure:"HTTP_PORT"`
	SecretKey          string             `mapstructure:"SECRET_KEY"`
	MongoDtabaseConfig MongoDtabaseConfig `mapstructure:",squash"`
}

type MongoDtabaseConfig struct {
	URI  string `mapstructure:"MONGO_DATABASE_URI"`
	Name string `mapstructure:"MONGO_DATABASE_NAME"`
}

func Load() (Config, error) {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")

	viper.ReadInConfig()
	// if err != nil {
	// }
	viper.BindEnv("HTTP_URI")
	viper.BindEnv("PORT")
	viper.BindEnv("SECRET_KEY")
	viper.BindEnv("MONGO_DATABASE_URI")
	viper.BindEnv("MONGO_DATABASE_NAME")

	var c Config
	err := viper.Unmarshal(&c)

	return c, err
}
