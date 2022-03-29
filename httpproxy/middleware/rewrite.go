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

////进行tls握手
//func dialTLS(network, addr string) (net.Conn, error) {
//	conn, err := net.Dial(network, addr)
//	if err != nil {
//		return nil, err
//	}
//
//	host, _, err := net.SplitHostPort(addr)
//	if err != nil {
//		return nil, err
//	}
//	cfg := &tls.Config{ServerName: host}
//
//	tlsConn := tls.Client(conn, cfg)
//	if err := tlsConn.Handshake(); err != nil {
//		conn.Close()
//		return nil, err
//	}
//
//	cs := tlsConn.ConnectionState()
//	cert := cs.PeerCertificates[0]
//
//	// Verify here
//	cert.VerifyHostname(host)
//	log.Println(cert.Subject)
//
//	return tlsConn, nil
//}
func customTransport(to *url.URL) http.RoundTripper {
	var dialer = &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
		//DualStack: true,
	}

	var transport http.RoundTripper = &http.Transport{
		//Proxy: nil, // http.ProxyFromEnvironment,
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			return dialer.DialContext(ctx, network, addr)
		},
		MaxIdleConns:          8,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
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
func NewUrlRewriteMiddleWare(from, to *url.URL) *RewriteUrlMiddleWare {

	var transport = customTransport(to)

	revProxy := httputil.NewSingleHostReverseProxy(to)
	//设置连接池
	revProxy.Transport = transport
	//异常回调
	revProxy.ErrorHandler = httpcallback.OnError
	//revProxy.Director = func(w *http.Request) {
	//
	//}

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
var hopHeaders = []string{
	"Connection",
	"Proxy-Connection", // non-standard but still sent by libcurl and rejected by e.g. google
	"Keep-Alive",
	"Proxy-Authenticate",
	"Proxy-Authorization",
	"Te",      // canonicalized version of "TE"
	"Trailer", // not Trailers per URL above; http://www.rfc-editor.org/errata_search.php?eid=4522
	"Transfer-Encoding",
	"Upgrade",
}

func removeHeaders(header http.Header) {
	// Remove hop-by-hop headers listed in the "Connection" header.
	if c := header.Get("Connection"); c != "" {
		for _, f := range strings.Split(c, ",") {
			if f = strings.TrimSpace(f); f != "" {
				header.Del(f)
			}
		}
	}

	// Remove hop-by-hop headers
	for _, h := range hopHeaders {
		if header.Get(h) != "" {
			header.Del(h)
		}
	}
}
func addXForwardedForHeader(req *http.Request) {
	if clientIP, _, err := net.SplitHostPort(req.RemoteAddr); err == nil {
		// If we aren't the first proxy retain prior
		// X-Forwarded-For information as a comma+space
		// separated list and fold multiple headers into one.
		if prior, ok := req.Header["X-Forwarded-For"]; ok {
			clientIP = strings.Join(prior, ", ") + ", " + clientIP
		}
		req.Header.Set("X-Forwarded-For", clientIP)
	}
}
func (p *RewriteUrlMiddleWare) ProxyHTTPS(rw http.ResponseWriter, req *http.Request) {
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
	var (
		raw, old, now = strutil.String2Bytes(&r.URL.Path), strutil.String2Bytes(&ws.From.Path), strutil.String2Bytes(&ws.To.Path)
	)
	r.URL.Path = strutil.Byte2String(strutil.UnsafeReplaceBegin(raw, old, now))
	addXForwardedForHeader(r)
	removeHeaders(r.Header)
	if r.Method == http.MethodConnect {
		ws.ProxyHTTPS(w, r)
		return
	}
	// connect

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
