package httpproxy

import (
	"github.com/spf13/viper"
	"log"
	"myproxyHttp/httpproxy/middleware"
	"net/http"
	"net/url"
)

func Simple(from, to *url.URL) {

	var s = middleware.NewUrlRewriteMiddleWare(from, to)
	http.HandleFunc(from.Path, func(w http.ResponseWriter, r *http.Request) {
		//log.Printf("detail %v\n", r.Header)
		s.ServeHTTP(w, r)
	})
	//log.Printf("ping test ")
	err := s.LiveCheck()
	log.Printf("information = > %v  ", err)
	http.ListenAndServe(from.Host, nil)

}

func Simple1() {

	var s1 = []string{
		"http://localhost:8081/",
		"http://localhost:8083/",
		"http://localhost:8082/",
		"http://localhost:8084/",
		"http://localhost:10000",
	}
	var s, err = middleware.NewLbProxyHandler("http://localhost:8080/", s1)
	if err != nil {
		panic(err)
	}
	//s.UseConsistHash()
	log.Printf("begin it\n")
	go middleware.ClearThread.Process()
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		s.ServeHTTP(writer, request)
	})
	http.ListenAndServe(":8080", nil)

}

func _reg(f string, to []interface{}) {
	go func() {
		if len(to) == 0 {
			return
		}
		var s = make([]string, len(to))
		for i, _ := range to {
			s[i] = to[i].(string)
		}
		handler, err := middleware.NewLbProxyHandler(f, s)
		if err != nil {
			log.Printf("error %+v")
			return
		}
		var fhttp = http.NewServeMux()
		fhttp.HandleFunc(handler.From.Path, func(w http.ResponseWriter, r *http.Request) {
			handler.ServeHTTP(w, r)
		})
		log.Printf("init http proxy from=[%s], to=[%v] \n", f, to)
		go middleware.ClearThread.Process()
		http.ListenAndServe(handler.From.Host, fhttp)

	}()
}
func WithConfig() {
	_array := viper.Get("proxy.http")

	arr, isArr := _array.([]interface{})
	if !isArr {
		return
	}

	for _, o := range arr {
		log.Printf("%T\n", o)
		mp, ismp := o.(map[interface{}]interface{})
		if !ismp {

			panic("proxy.http is not array")
		}
		if s, ok := mp["nodes"]; !ok {
			log.Printf("v = %v", s)
		}
		m, ok := mp["nodes"].(map[interface{}]interface{})
		if !ok {
			panic("nodes is not object")
		}
		//is mp
		s := m["from"].(string)
		t := m["to"].([]interface{})
		go _reg(s, t)
	}

}
