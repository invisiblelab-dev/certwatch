package commands

import (
	"github.com/spf13/viper"
)

func Parse() {
	viper.SetConfigName("certwatch")
	viper.SetConfigType("yml")
	viper.AddConfigPath(".")

	viper.ReadInConfig()
}
