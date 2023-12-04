package conf

import (
	"github.com/spf13/viper"
	"log"
)

type Configuration struct {
	Url string
}

func New() Configuration {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}

	return Configuration{
		Url: viper.Get("URL").(string),
	}
}
