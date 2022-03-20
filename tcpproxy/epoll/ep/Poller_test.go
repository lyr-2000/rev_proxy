package ep

import (
	"log"
	"syscall"
	"testing"
)

func TestListen(t *testing.T) {
	var p = Poller{}
	p.OnConnOpen = func(epfd EpollFd, conn SockFd) {
		log.Printf("conn open %+v\n", conn)
	}
	p.OnFdRemoved = func(epfd EpollFd, conn SockFd) {
		conn.Close()
	}
	p.OnMsgReceive = func(fd EpollFd, conn SockFd, eventCode uint32) {
		var buf = make([]byte, 4096)
		n, err := conn.Read(buf)
		log.Printf("event code is {%v}", eventCode)
		if err != nil {
			log.Printf("read error %+v", err)
		}
		if n == 0 {
			log.Printf("receive is zero ,close fd")
			EpollRemove(fd, conn, &p)
			return
		}
		conn.Write([]byte("hello world\n"))
		conn.Write(buf)
	}
	var q [100]syscall.EpollEvent
	_, err := Listen(&p, "", 8080, 10, q[:])
	if err != nil {
		panic(err)
	}

}
