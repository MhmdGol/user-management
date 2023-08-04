package config

import "github.com/spf13/viper"

type Config struct {
	HttpURI            string             `mapstructure:"HTTP_URI"`
	Port               string             `mapstructure:"HTTP_PORT"`
	RsaPair            RsaPair            `mapstructure:",squash"`
	MongoDtabaseConfig MongoDtabaseConfig `mapstructure:",squash"`
}

type RsaPair struct {
	SecretKeyPath string `mapstructure:"SECRET_KEY_PATH"`
	PublicKeyPath string `mapstructure:"PUBLIC_KEY_PATH"`
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
	viper.BindEnv("HTTP_PORT")
	viper.BindEnv("SECRET_KEY_PATH")
	viper.BindEnv("PUBLIC_KEY_PATH")
	viper.BindEnv("MONGO_DATABASE_URI")
	viper.BindEnv("MONGO_DATABASE_NAME")

	var c Config
	err := viper.Unmarshal(&c)

	return c, err
}
