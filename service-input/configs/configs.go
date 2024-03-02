package configs

import (
	"github.com/spf13/viper"
)

func LoadConfig(path string) error {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	viper.AutomaticEnv()
	return nil
}
