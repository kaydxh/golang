package viper

import (
	"github.com/ory/viper"
)

func GetViper(configFile string, subKey string) *viper.Viper {

	if configFile == "" {
		v := viper.GetViper()
		if v == nil {
			return nil
		}

		return v
	}

	viper.SetConfigFile(configFile)
	err := viper.ReadInConfig()
	if err != nil {
		return nil
	}

	if subKey == "" {
		return viper.GetViper()
	}

	return viper.Sub(subKey)
}
