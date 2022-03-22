package middleware

import (
	"log"
	"sync"
	"testing"
	"time"
)

func Test_mutex0(t *testing.T) {

	var s sync.RWMutex
	s.Lock()

	log.Printf("获取 锁成功")
	s.Unlock()
	s.RLock()
	defer s.RUnlock()
	log.Printf("获取 读锁成功")

	for {
		time.Sleep(time.Second * 10)
	}

}
