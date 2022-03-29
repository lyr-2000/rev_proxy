package middleware

import (
	"context"
	"crypto/tls"
	"io"
	"log"
	"myproxyHttp/httpproxy/httpcallback"
	"myproxyHttp/utils/strutil"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
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

func customTransport(to *url.URL) http.RoundTripper {
	var dialer = &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
		//DualStack: true,
	}

	var transport http.RoundTripper = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			//network == tcp and  addr == 443
			//log.Printf("dial context %v, %v", network, addr)
			return dialer.DialContext(ctx, network, addr)
		},
		MaxIdleConns:          8,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 10 * time.Second,
		DisableKeepAlives:     false,
		//DialTLSContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
		//	conn, err := dialTLS(network, addr)
		//	return conn, err
		//},
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, //跳过 https验证

		//DisableCompression: true,
		//dialt,

	}
	//if to != nil && to.Scheme == "https" {
	//	//https wetting
	//	//transport.
	//}
	return transport
}

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}
func joinURLPath(a, b *url.URL) (path, rawpath string) {
	if a.RawPath == "" && b.RawPath == "" {
		return singleJoiningSlash(a.Path, b.Path), ""
	}
	// Same as singleJoiningSlash, but uses EscapedPath to determine
	// whether a slash should be added
	apath := a.EscapedPath()
	bpath := b.EscapedPath()

	aslash := strings.HasSuffix(apath, "/")
	bslash := strings.HasPrefix(bpath, "/")

	switch {
	case aslash && bslash:
		return a.Path + b.Path[1:], apath + bpath[1:]
	case !aslash && !bslash:
		return a.Path + "/" + b.Path, apath + "/" + bpath
	}
	return a.Path + b.Path, apath + bpath
}
func NewUrlRewriteMiddleWare(from, to *url.URL) *RewriteUrlMiddleWare {

	var transport = customTransport(to)
	var (
		old, now = strutil.String2Bytes(&from.Path), strutil.String2Bytes(&to.Path)
	)
	targetQuery := to.RawQuery
	director := func(req *http.Request) {
		//禁用重定向

		req.URL.Scheme = to.Scheme
		req.URL.Host = to.Host

		req.URL.Path, req.URL.RawPath = joinURLPath(to, req.URL)
		//replace the url path  :  localhost:8080/api/xxx =>  www.baidu.com/xxx
		raw := strutil.String2Bytes(&req.URL.Path)
		req.URL.Path = strutil.Byte2String(strutil.UnsafeReplaceBegin(raw, old, now))
		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}
		if _, ok := req.Header["User-Agent"]; !ok {
			// explicitly disable User-Agent so it's not set to default value
			req.Header.Set("User-Agent", "")
		}
		//修改 localhost 记录
		req.Host = to.Host

	}

	revProxy := &httputil.ReverseProxy{
		Director:  director,
		Transport: transport,
	}
	//设置连接池
	//revProxy.Transport = transport
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

// Hop-by-hop headers. These are removed when sent to the backend.
// http://www.w3.org/Protocols/rfc2616/rfc2616-sec13.html
//var hopHeaders = []string{
//	"Connection",
//	"Proxy-Connection", // non-standard but still sent by libcurl and rejected by e.g. google
//	"Keep-Alive",
//	"Proxy-Authenticate",
//	"Proxy-Authorization",
//	"Te",      // canonicalized version of "TE"
//	"Trailer", // not Trailers per URL above; http://www.rfc-editor.org/errata_search.php?eid=4522
//	"Transfer-Encoding",
//	"Upgrade",
//}
//
//func removeHeaders(header http.Header) {
//	// Remove hop-by-hop headers listed in the "Connection" header.
//	if c := header.Get("Connection"); c != "" {
//		for _, f := range strings.Split(c, ",") {
//			if f = strings.TrimSpace(f); f != "" {
//				header.Del(f)
//			}
//		}
//	}
//
//	// Remove hop-by-hop headers
//	for _, h := range hopHeaders {
//		if header.Get(h) != "" {
//			header.Del(h)
//		}
//	}
//}
//func addXForwardedForHeader(req *http.Request) {
//	if clientIP, _, err := net.SplitHostPort(req.RemoteAddr); err == nil {
//		// If we aren't the first proxy retain prior
//		// X-Forwarded-For information as a comma+space
//		// separated list and fold multiple headers into one.
//		if prior, ok := req.Header["X-Forwarded-For"]; ok {
//			clientIP = strings.Join(prior, ", ") + ", " + clientIP
//		}
//		req.Header.Set("X-Forwarded-For", clientIP)
//	}
//}
func (p *RewriteUrlMiddleWare) ProxyHTTPS(rw http.ResponseWriter, req *http.Request) {
	req.Host = p.To.Host
	hij, ok := rw.(http.Hijacker)
	if !ok {
		log.Printf("http server does not support hijacker")
		return
	}

	clientConn, _, err := hij.Hijack()
	if err != nil {
		log.Printf("http: proxy error: %v", err)
		return
	}

	proxyConn, err := net.Dial("tcp", req.URL.Host)
	if err != nil {
		log.Printf("http: proxy error: %v", err)
		return
	}

	// The returned net.Conn may have read or write deadlines
	// already set, depending on the configuration of the
	// Server, to set or clear those deadlines as needed
	// we set timeout to 5 minutes
	deadline := time.Now()
	//if p.Timeout == 0 {
	deadline = deadline.Add(time.Minute * 2)

	err = clientConn.SetDeadline(deadline)
	if err != nil {
		log.Printf("http: proxy error: %v", err)
		return
	}

	err = proxyConn.SetDeadline(deadline)
	if err != nil {
		log.Printf("http: proxy error: %v", err)
		return
	}

	_, err = clientConn.Write([]byte("HTTP/1.0 200 OK\r\n\r\n"))
	if err != nil {
		log.Printf("http: proxy error: %v", err)
		return
	}

	go func() {
		io.Copy(clientConn, proxyConn)
		clientConn.Close()
		proxyConn.Close()
	}()

	io.Copy(proxyConn, clientConn)
	proxyConn.Close()
	clientConn.Close()
}
func (ws *RewriteUrlMiddleWare) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodConnect {
		ws.ProxyHTTPS(w, r)
		return
	}

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
