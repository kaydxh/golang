/*
 *Copyright (c) 2022, kaydxh
 *
 *Permission is hereby granted, free of charge, to any person obtaining a copy
 *of this software and associated documentation files (the "Software"), to deal
 *in the Software without restriction, including without limitation the rights
 *to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 *copies of the Software, and to permit persons to whom the Software is
 *furnished to do so, subject to the following conditions:
 *
 *The above copyright notice and this permission notice shall be included in all
 *copies or substantial portions of the Software.
 *
 *THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 *IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 *FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 *AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 *LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 *OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 *SOFTWARE.
 */
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

func GetNewViper(configFile string, subKeys string) *viper.Viper {
	if configFile == "" {
		return nil
	}

	viperNew := viper.New()
	viperNew.SetConfigFile(configFile)
	err := viperNew.ReadInConfig()
	if err != nil {
		return nil
	}

	if subKeys == "" {
		return viperNew
	}

	v := viperNew
	keys := strings.Split(subKeys, ".")

	for _, k := range keys {
		v = v.Sub(k)
		if v == nil {
			return nil
		}
	}

	return v
}
