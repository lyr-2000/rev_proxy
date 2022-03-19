package httpproxy

import (
	"errors"
	"fmt"
	"log"
	"myproxyHttp/utils/reflectutil"
	"myproxyHttp/utils/strutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

var (
	portRouter = map[string]*http.ServeMux{}
)

func DoneHttpWait() {
	for port, _ := range portRouter {
		var fakeHttp = portRouter[port]
		if fakeHttp != nil {
			go func(port string, fakeHttp *http.ServeMux) {
				log.Printf("start proxy [localhost:%v] ", port)

				err := http.ListenAndServe(fmt.Sprintf(":%s", port), fakeHttp)

				if err != nil {
					log.Fatalf("listen error %+v\n", err)
				}
			}(port, fakeHttp)
		}
		delete(portRouter, port)
	}
}

func standardUrl(from *url.URL) {
	if from.Path == "" {
		from.Path = "/"
	} else if from.Path[len(from.Path)-1] != '/' {
		from.Path = from.Path + "/"
	}

}
func rewritePath(r *http.Request, from *url.URL, to *url.URL) {
	// /ide/index.html            /ide/          /abc/
	// /abc/index.html
	w := strings.Replace(r.URL.Path, from.Path, to.Path, 1)
	r.URL.Path = w
	//fmt.Printf("path %v, from %v,to %v  w %v\n", r.URL.Path, from.Path, to.Path, w)
}
func registerProxyConf(from, to string) error {
	fromUrl, err := url.Parse(from)
	if err != nil {
		return err
	}
	toUrl, err := url.Parse(to)
	if err != nil {
		return err
	}
	standardUrl(fromUrl)
	standardUrl(toUrl)
	path := fromUrl.Path
	log.Printf("[%v, %v]\n", fromUrl, toUrl)
	//if path == "" {
	//	path = "/"
	//} else {
	//	if path[len(path)-1] != '/' {
	//		path = path + "/"
	//	}
	//	log.Printf("path = > [%v]\n", path)
	//}
	// forward remote => 代理到 toUrl
	proxy := httputil.NewSingleHostReverseProxy(toUrl)
	//fakeHttp := http.ServeMux{}
	var fakeHttp = portRouter[fromUrl.Port()]
	if fakeHttp == nil {
		fakeHttp = http.NewServeMux()
		portRouter[fromUrl.Port()] = fakeHttp
	}

	// to 就是转发到 本机的端口
	fakeHttp.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		rewritePath(r, fromUrl, toUrl)
		proxy.ServeHTTP(w, r)
	})
	log.Printf("config  proxy [%v]=> [%v]\n", fromUrl, toUrl)
	//go func() {
	//	log.Printf("start proxy [%v]=> [%v]", fromUrl, toUrl)
	//
	//	err := http.ListenAndServe(fmt.Sprintf(":%v", fromUrl.Port()), fakeHttp)
	//
	//	if err != nil {
	//		log.Fatalf("listen error %+v\n", err)
	//	}
	//}()

	return nil

}

func ParseConfigMapDefault(mp map[string]interface{}) error {
	if mp == nil {
		log.Println("无http配置")
		return nil
	}
	log.Printf("config map %#v", strutil.ToJSON(mp))
	var proxy = mp["proxy"].(map[string]interface{})
	if httpConf, ok := proxy["http"]; ok {
		if !reflectutil.IsArrayOrSlice(httpConf) {
			// 如果不是数组
			return errors.New("配置不正确，proxy.http不是数组")
		}
		arr := httpConf.([]interface{})
		for i, _ := range arr {
			h := arr[i].(map[string]interface{})
			err := registerProxyConf(h["from"].(string), h["to"].(string))
			if err != nil {
				log.Printf("%+v", err)
			}
		}
		// proxy parse
	}

	return nil

}
