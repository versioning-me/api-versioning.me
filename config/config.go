package config

import (
	"github.com/spf13/viper"
)

var (
	Config *viper.Viper
	Env    string
)

func Init(env string) error {
	var err error
	Config = viper.New()
	Config.SetConfigType("yaml")
	Config.SetConfigName(env)
	Config.AddConfigPath("../config/")
	Config.AddConfigPath("config/")
	err = Config.ReadInConfig()
	if err != nil {
		return err
	}
	Env = Config.Get("environment").(string)
	return nil
}
