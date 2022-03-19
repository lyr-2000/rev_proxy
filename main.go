package main

import (
	"flag"
	yaml "github.com/goccy/go-yaml"
	"io/ioutil"
	"log"
	"myproxyHttp/httpproxy"
	"myproxyHttp/tcpproxy"
	"myproxyHttp/utils/fileutil"
	"time"
)

func main() {
	confPath := *flag.String("conf", "conf", "configFile path")
	//fmt.Println("hello world")
	flag.Parse()
	_ = fileutil.MkdirAll(confPath)
	list := fileutil.ListAllFilePathInDir(confPath)
	for _, v := range list {
		bytes, err := ioutil.ReadFile(v)
		if err != nil {
			log.Printf("io error %+v\n", err)
			continue
		}
		var s map[string]interface{}
		if err := yaml.Unmarshal(bytes, &s); err != nil {
			//
			log.Printf("error info %+v", s)
			continue
		}
		// 注册http路由配置
		err = httpproxy.ParseConfigMapDefault(s)
		if err != nil {
			log.Printf("error info %+v", err)
		}
		//all bytes
		//注册tcp配置
		err = tcpproxy.ParseConfigTcp(s)
		if err != nil {
			log.Printf("error info %+v", err)
		}

	}
	//正式开始代理 http
	httpproxy.DoneHttpWait()

	tcpproxy.DoneFinal()

	//开启tcp代理
	for {
		time.Sleep(time.Hour)
	}
}
