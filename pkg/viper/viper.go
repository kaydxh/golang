package viper

import (
	"strings"

	"github.com/spf13/viper"
)

//keys delim with dot(.), eg: "database.mysql"
func GetViper(configFile string, subKeys string) *viper.Viper {
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

	if subKeys == "" {
		return viper.GetViper()
	}

	v := viper.GetViper()
	keys := strings.Split(subKeys, ".")

	for _, k := range keys {
		v = v.Sub(k)
		if v == nil {
			return nil
		}
	}

	return v
}
