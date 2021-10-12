package env

import (
	"github.com/spf13/viper"
	"github.com/symmetric-project/symphony/utils"
)

type Config struct {
	DATABASE_URL string `mapstructure:"DATABASE_URL"`
}

var CONFIG Config

func init() {
	viper.SetConfigFile(".env")

	if err := viper.ReadInConfig(); err != nil {
		utils.StacktraceErrorAndExit(err)
	}
	if err := viper.Unmarshal(&CONFIG); err != nil {
		utils.StacktraceErrorAndExit(err)
	}
}
