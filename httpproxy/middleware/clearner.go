package middleware

import (
	"log"
	"sync"
	"time"
)

//用来清理掉部分 down掉的主机
type clearThread_ struct {
	Channel  chan *LbProxyHandler
	isActive bool
	Queue    []*LbProxyHandler
	sync.Mutex
}

func makeclearThread() *clearThread_ {
	var x = clearThread_{}
	x.Channel = make(chan *LbProxyHandler, 10)
	return &x
}

var ClearThread = makeclearThread()

func (u *clearThread_) Process() {

	if u.isActive {
		return
	} else {
		u.Lock()
		if u.isActive {
			u.Unlock()
			return
		}
		u.isActive = true
		u.Unlock()
	}

	//var isActive = true
	for {

		select {
		case proxyHandler := <-u.Channel:
			//isActive = true
			log.Printf("host down, adjust prehost = %v\n", proxyHandler.AliveHost)
			adjustProxyObject(proxyHandler)
			log.Printf("adjusted alive host = %v\n", proxyHandler.AliveHost)
		case <-time.After(time.Minute * 3):
			//上次没有通知，这次也没有

		}
	}
}

func adjustProxyObject(u *LbProxyHandler) {
	if u == nil {
		return
	}

	u.clearDiedProxy()

}
