package epoll

import (
	"log"
	"myproxyHttp/consts"
	"myproxyHttp/tcpproxy/epoll/ep"
	"net/url"
	"strconv"
	"syscall"
)

type EpHandler struct {
	mp map[ep.SockFd]ep.SockFd
}

func (handler *EpHandler) Serve(from, to *url.URL) {
	go handler.Serve0(from, to)
}
func (handler *EpHandler) Serve0(from, to *url.URL) {
	hostname := from.Hostname()
	iport, _ := strconv.Atoi(from.Port())
	handler.mp = make(map[ep.SockFd]ep.SockFd, consts.TcpSockFDMapSize)

	var poller = ep.Poller{}
	//被移除epoll监听的时候,把网络连接关闭
	poller.OnFdRemoved = func(epfd ep.EpollFd, conn ep.SockFd) {
		var value = handler.mp[conn]
		//因为目前是 单线程的事件循环，不需要加锁，多个协程就要考虑加锁
		delete(handler.mp, conn)
		//delete  conn => value
		// value => conn
		if handler.mp[value] == conn {
			ep.EpollRemove(epfd, value, &poller)
			//delete(handler.mp, value)
		}
		//epoll不会自动关闭 ，要我自己手动关闭 比较保险
		conn.Close()
		log.Printf("close fd [%v, %v]\n", conn, value)
	}
	//连接监听时候的回调
	poller.OnConnOpen = func(epfd ep.EpollFd, conn ep.SockFd) {
		sockfd, c, err := ep.Open("tcp", to.Host)
		if err != nil || sockfd <= 0 || c == nil {
			ep.EpollRemove(epfd, conn, &poller)
			return
		}
		ep.EpollCtl(epfd, sockfd, syscall.EPOLLIN|syscall.EPOLLPRI)
		//创建连接
		handler.mp[sockfd] = conn
		handler.mp[conn] = sockfd
	}
	var buf = make([]byte, consts.TcpBufSize)
	poller.OnMsgReceive = func(epfd ep.EpollFd, conn ep.SockFd, eventCode uint32) {
		nbuf, err := conn.Read(buf)
		if err != nil {
			ep.EpollRemove(epfd, conn, &poller)
			return
		}
		if nbuf <= 0 {
			ep.EpollRemove(epfd, conn, &poller)
			return
		}
		nw, err := handler.mp[conn].Write(buf[:nbuf])
		if err != nil {
			log.Printf("write error %v", err)
			ep.EpollRemove(epfd, conn, &poller)
			return
		}
		if nw <= 0 {
			ep.EpollRemove(epfd, conn, &poller)
			return
		}

	}

	var queue = make([]syscall.EpollEvent, consts.TcpEpollEventQueueSize)
	//go func() {
	_, err := ep.Listen(&poller, hostname, iport, consts.TcpBackLog, queue[:])
	if err != nil {
		log.Printf("epoll error %v", err)
	}
	//}()

}
