package middleware

import (
	"net/http"
)

type ProxyHandler struct {
	callback func(w http.ResponseWriter, r *http.Request)
	//logic    func(w http.ResponseWriter, r *http.Request, next func(w http.ResponseWriter, r *http.Request))
}

func (f *ProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f.callback(w, r)
}

func (f *ProxyHandler) Wrap(f1 func(w http.ResponseWriter, r *http.Request, next func(w http.ResponseWriter, r *http.Request))) *ProxyHandler {
	pre := f.callback
	f.callback = func(w http.ResponseWriter, r *http.Request) {
		f1(w, r, pre)
	}
	return f
}

/*
func (f *ProxyHandler) Proxyexample() {
	var s ProxyHandler
	s.Wrap(func(w http.ResponseWriter, r *http.Request, next func(w http.ResponseWriter, r *http.Request)) {
		if r.URL.Path == "/admin" {
			w.Write([]byte("大胆"))
			return
		}

		next(w, r)
	})
}


*/
