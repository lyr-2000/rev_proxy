package middleware

import (
	"myproxyHttp/httpproxy/loadbalance/consistenhash"
	"myproxyHttp/utils/serverutil"
	"net/http"
)

//import "myproxyHttp/httpproxy/loadbalance/consistenhash"

type ConsistHashHandler struct {
	Ring *consistenhash.Ring
}

func (x *ConsistHashHandler) Each(f func(i int, u *RewriteUrlMiddleWare)) {
	for i, v := range x.Ring.Nodes {
		f(i, v.Value.(*RewriteUrlMiddleWare))
	}
}

func (x *ConsistHashHandler) RemoveProxy(u *RewriteUrlMiddleWare) error {
	err := x.Ring.RemoveNode(u.To.Host)
	return err
}
func NewConsistHashHandler(Handler []*RewriteUrlMiddleWare) *ConsistHashHandler {
	ring := consistenhash.NewRing()
	for _, v := range Handler {
		ring.AddNode(v.To.Host, v.Proxy)
	}

	return &ConsistHashHandler{ring}
}
func (ring *ConsistHashHandler) Get(w http.ResponseWriter, r *http.Request) http.Handler {

	ok := ring.Ring.Get(serverutil.ClientIP(r))
	if ok != nil {
		handler, ishandler := ok.(http.Handler)
		if ishandler {
			return handler
		}
	}
	return nil
}

func (r *ConsistHashHandler) Has(id string) bool {
	return r.Ring.Has(id)
}
