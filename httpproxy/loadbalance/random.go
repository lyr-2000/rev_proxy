package loadbalance

import (
	"math/rand"
	"time"
)

/*
1. 随机负载
2. 轮询负载
3. 加权负载
4. 一致性hash负载

*/
type SimpleRandomLb struct {
}

func (u *SimpleRandomLb) NextIndex(size int) int {
	return rand.Intn(size)
}

func init() {
	rand.Seed(time.Now().UnixMilli())
}
