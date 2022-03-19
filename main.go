package main

import (
	"fmt"
	yaml "github.com/goccy/go-yaml"
	"io/ioutil"
	"log"
	"myproxyHttp/httpproxy"
	"myproxyHttp/utils/fileutil"
	"time"
)

func main() {
	fmt.Println("hello world")
	list := fileutil.ListAllFilePathInDir("./conf")
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
		// travel
		err = httpproxy.ParseConfigMapDefault(s)
		if err != nil {
			log.Printf("error info %+v", err)
		}
		//all bytes

	}
	//正式开始代理 http
	httpproxy.DoneHttpWait()
	for {
		time.Sleep(time.Hour)
	}
}
