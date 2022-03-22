package main

import (
	"flag"
	"log"
	"myproxyHttp/config"
	"myproxyHttp/httpproxy"
	"myproxyHttp/tcpproxy"
	"time"
)

func main() {
	cconf := *flag.String("conf", "conf", "configFile path")
	//fmt.Println("hello world")
	flag.Parse()
	//_ = fileutil.MkdirAll(confPath)
	//list := fileutil.ListAllFilePathInDir(confPath)
	//config.Load(confPath)

	//from, _ := url.Parse("http://localhost:8080/")
	//to, _ := url.Parse("http://localhost:10004/")
	config.Load(cconf)
	httpproxy.WithConfig()
	tcpproxy.WithConfig()
	//httpproxy.Simple1()
	//开启tcp代理
	log.Printf("init done\n")
	for {
		time.Sleep(time.Hour)
	}
}
