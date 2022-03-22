package loadbalance

import (
	"math/rand"
	"net/http"
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

func NewSimpleRandomLb() *SimpleRandomLb {
	return &SimpleRandomLb{}
}
func (u *SimpleRandomLb) NextIndex(w http.ResponseWriter, r *http.Request, aliveHost int) int {
	if aliveHost <= 0 {
		return 0
	}
	return rand.Intn(aliveHost)
}
func (u *SimpleRandomLb) RemoveHost(key string) {

}

func init() {
	rand.Seed(time.Now().UnixMilli())
}

type LbObject interface {
	NextIndex(w http.ResponseWriter, r *http.Request, aliveHost int) int

	RemoveHost(key string)
}
