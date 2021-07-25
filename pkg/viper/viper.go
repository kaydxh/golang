package viper

import (
	"fmt"

	"github.com/ory/viper"
)

func GetViper(configFile string, subKey string) (*viper.Viper, error) {

	if configFile == "" {
		v := viper.GetViper()
		if v == nil {
			return nil, fmt.Errorf("failed to get viper without config file")
		}

		return v, nil
	}

	viper.SetConfigFile(configFile)
	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to read config err: %v", err)
	}

	if subKey == "" {
		return viper.GetViper(), nil
	}

	return viper.Sub(subKey), nil
}
