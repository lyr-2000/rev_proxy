package config

import (
	"fmt"
	"github.com/spf13/viper"
	"myproxyHttp/utils/fileutil"
	"strings"
)

func Load(confPath string) {
	dir := fileutil.ListAllFilePathInDir(confPath)
	if len(dir) <= 0 {
		panic(fmt.Sprintf("no config file in %s\n", confPath))
	}
	for _, v := range dir {
		if strings.HasSuffix(v, ".yml") ||
			strings.HasSuffix(v, ".yaml") {
			viper.SetConfigFile(v)
		}
	}
	//viper.SetConfigFile(dir[0])
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
