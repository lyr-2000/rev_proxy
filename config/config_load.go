package config

import (
	"fmt"
	"github.com/spf13/viper"
)

func Load(confPath string) {

	viper.SetConfigName(confPath)
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal read config file %+v\n", err))
	}

}

//设置常量
func InitConfigVars() {

}

//清除变量
func ClearConfigVars() {

}
