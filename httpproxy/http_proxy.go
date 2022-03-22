package httpproxy

import (
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
