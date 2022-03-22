package middleware

import (
	"context"
	"log"
	"myproxyHttp/httpproxy/httpcallback"
	"myproxyHttp/utils/strutil"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"
)

type RewriteUrlMiddleWare struct {
	From, To    *url.URL
	Proxy       *httputil.ReverseProxy
	FailDialCnt int64
	Down        bool

	sync.Mutex
}

func (u *RewriteUrlMiddleWare) IsDeath() bool {
	return u.Down
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

func customTransport() http.RoundTripper {
	var dialer = &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
		//DualStack: true,
	}

	var transport http.RoundTripper = &http.Transport{
		Proxy: nil, // http.ProxyFromEnvironment,
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			return dialer.DialContext(ctx, network, addr)
		},
		MaxIdleConns:          8,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		//DisableCompression: true,
	}
	return transport
}
func NewUrlRewriteMiddleWare(from, to *url.URL) *RewriteUrlMiddleWare {

	var transport = customTransport()

	revProxy := httputil.NewSingleHostReverseProxy(to)
	//设置连接池
	revProxy.Transport = transport
	//异常回调
	revProxy.ErrorHandler = httpcallback.OnError

	//修改回调
	revProxy.ModifyResponse = httpcallback.ModifyResponse

	//设置异常回调
	return &RewriteUrlMiddleWare{
		From:  from,
		To:    to,
		Proxy: revProxy,
	}
}

type HttpHandler func(w http.ResponseWriter, r *http.Request)

func (ws *RewriteUrlMiddleWare) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		raw, old, now = strutil.String2Bytes(&r.URL.Path), strutil.String2Bytes(&ws.From.Path), strutil.String2Bytes(&ws.To.Path)
	)
	r.URL.Path = strutil.Byte2String(strutil.UnsafeReplaceBegin(raw, old, now))
	//直接代理转发
	ws.Proxy.ServeHTTP(w, r)
}

func (ws *RewriteUrlMiddleWare) Ping() (*http.Response, error) {
	request, err := http.NewRequest("GET", ws.To.String(), nil)
	if err != nil {
		return nil, err
	}
	return ws.Proxy.Transport.RoundTrip(request)
}

func retry(ws *RewriteUrlMiddleWare) {
	if ws == nil {
		return
	}
	var sleepTime = 1
	for {
		if sleepTime > 10000 {
			sleepTime = 0
		}
		time.Sleep(time.Second * 3)
		_, err := ws.Ping()
		if err == nil {
			ws.Lock()
			// change  status
			ws.Down = false
			ws.FailDialCnt = 0
			ws.Unlock()
			break
		} else {
			// is not nil
			time.Sleep(time.Second * time.Duration(sleepTime))
			sleepTime += 10
			log.Printf("host %v is down \n", ws.To)

		}

	}

}
func (ws *RewriteUrlMiddleWare) LiveCheck() error {
	//if ws.Down {
	//	return errors.New("node is down")
	//}
	_, err := ws.Ping()

	if err != nil {
		log.Printf("%v ,", err)
		ws.FailDialCnt++
		ws.Down = true
	} else {
		ws.Down = false
		ws.FailDialCnt = 0
	}
	// err != nil {
	return err

}
