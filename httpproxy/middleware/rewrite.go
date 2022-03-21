package middleware

import (
	"myproxyHttp/utils/strutil"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type RewriteUrlMiddleWare struct {
	From, To *url.URL
	Proxy    *httputil.ReverseProxy
}

//func (w *RewriteUrlMiddleWare) HostKey() string {
//	return w.To.Host
//}

func standardUrl(from *url.URL) {
	if from.Path == "" {
		from.Path = "/"
	} else if from.Path[len(from.Path)-1] != '/' {
		from.Path = from.Path + "/"
	}

}

func NewUrlRewriteMiddleWare(from, to *url.URL) *RewriteUrlMiddleWare {
	//standardUrl(from)
	//standardUrl(to)

	return &RewriteUrlMiddleWare{
		From:  from,
		To:    to,
		Proxy: httputil.NewSingleHostReverseProxy(to),
	}
}

type HttpHandler func(w http.ResponseWriter, r *http.Request)

func (ws *RewriteUrlMiddleWare) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		raw, old, now = strutil.String2Bytes(r.URL.Path), strutil.String2Bytes(ws.From.Path), strutil.String2Bytes(ws.To.Path)
	)
	r.URL.Path = strutil.Byte2String(strutil.UnsafeReplaceBegin(raw, old, now))
	//直接代理转发
	ws.Proxy.ServeHTTP(w, r)
}
