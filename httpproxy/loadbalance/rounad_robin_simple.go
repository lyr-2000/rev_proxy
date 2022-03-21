package loadbalance

import "sync/atomic"

type RoundRobinSimple struct {
	Limit int32
	Index int32
}

func (lb *RoundRobinSimple) NextIndex() int32 {
	if lb.Limit <= 0 {
		return 0
	}
	res := atomic.LoadInt32(&lb.Index)
	//n-1 => 0
	if res >= lb.Limit-1 {
		//从0 开始
		atomic.CompareAndSwapInt32(&lb.Index, res, 0)
	} else {
		//incr
		atomic.CompareAndSwapInt32(&lb.Index, res, res+1)
	}
	return res % lb.Limit
}
