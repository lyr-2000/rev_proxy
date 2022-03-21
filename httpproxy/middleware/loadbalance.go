package middleware

import (
	"errors"
	"net/http"
	"net/url"
)

//负载均衡

type LbProxyHandler struct {
	From    *url.URL
	ToList  []*url.URL
	Handler []http.Handler
	//ProxyList []*httputil.ReverseProxy
}

func NewLbProxyHandler(fromU string, toUrlList []string) (*LbProxyHandler, error) {
	if len(toUrlList) == 0 {
		return nil, errors.New("no proxy target List")
	}
	from, err := url.Parse(fromU)
	if err != nil {
		return nil, err
	}
	standardUrl(from)

	var toList = make([]*url.URL, len(toUrlList), len(toUrlList))
	//var proxyList = make([]*httputil.ReverseProxy, len(toUrlList), len(toUrlList))
	var handler = make([]http.Handler, len(toUrlList), len(toUrlList))
	for i, v := range toUrlList {
		toList[i], err = url.Parse(v)
		if err != nil {
			return nil, err
		}
		//处理 后面的 /
		standardUrl(toList[i])
		handler[i] = NewUrlRewriteMiddleWare(from, toList[i])
		//proxyList[i] = httputil.NewSingleHostReverseProxy(toList[i])
	}
	return &LbProxyHandler{
		from, toList,
		handler,
	}, nil

}
