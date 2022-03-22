package loadbalance

import (
	"net/http"
	"sync/atomic"
)

//轮询的方式,用 cas的方式来做
type RoundRobinSimple struct {
	Index int32

	//sync.RWMutex
}

func NewRoundRobinSimple() *RoundRobinSimple {
	return &RoundRobinSimple{}
}
func (*RoundRobinSimple) RemoveHost(key string) {

}
func (lb *RoundRobinSimple) NextIndex(w http.ResponseWriter, r *http.Request, aliveHost int) int {
	if aliveHost <= 0 {
		return 0
	}
	//可能会remove Node
	//lb.RLock()
	//defer lb.RUnlock()

	limit := int32(aliveHost)
	if limit <= 0 {
		return 0
	}
	res := atomic.LoadInt32(&lb.Index)
	//n-1 => 0
	if res >= limit-1 {
		//从0 开始
		_ = atomic.CompareAndSwapInt32(&lb.Index, res, 0)
	} else {
		//incr
		_ = atomic.CompareAndSwapInt32(&lb.Index, res, res+1)
	}

	return int(res % limit)
}
