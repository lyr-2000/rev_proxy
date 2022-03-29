//go:build linux
// +build linux

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
	log.Printf("serve0 %v, %v", from, to)
	hostname := from.Hostname()
	iport, _ := strconv.Atoi(from.Port())
	handler.mp = make(map[ep.SockFd]ep.SockFd, consts.TcpSockFDMapSize)

	var poller = ep.Poller{}
	//被移除epoll监听的时候,把网络连接关闭
	poller.OnFdRemoved = func(epfd ep.EpollFd, conn ep.SockFd, _ *ep.Poller) {
		var value = handler.mp[conn]
		//因为目前是 单线程的事件循环，不需要加锁，多个协程就要考虑加锁
		delete(handler.mp, conn)
		//delete  conn => value
		// value => conn
		if handler.mp[value] == conn {
			_ = ep.EpollRemove(epfd, value, &poller)
			//delete(handler.mp, value)
		}
		//epoll不会自动关闭 ，要我自己手动关闭 比较保险
		_ = conn.Close()
		//log.Printf("close fd [%v, %v]\n", conn, value)
	}
	//连接监听时候的回调
	poller.OnConnOpen = func(epfd ep.EpollFd, conn ep.SockFd, _ *ep.Poller) {
		sockfd, c, err := ep.Open("tcp", to.Host)
		if err != nil || sockfd <= 0 || c == nil {
			_ = ep.EpollRemove(epfd, conn, &poller)
			return
		}
		sockfd.SetNonBlock()
		//客户端的话，最好监听 err 和 断开连接事件
		_ = ep.EpollCtl(epfd, sockfd, syscall.EPOLLIN|syscall.EPOLLPRI)
		//判断 ip地址校验之类的
		//{
		//	sockname, _ := ep.GetSockname(conn)
		//	ipv4, _ := ep.GetIpv4(sockname)
		//	log.Printf("ipv4 address %s", ipv4)
		//}
		//创建连接
		handler.mp[sockfd] = conn
		handler.mp[conn] = sockfd
	}
	var buf = make([]byte, consts.TcpBufSize)
	poller.OnSockFdActive = func(epfd ep.EpollFd, conn ep.SockFd, eventCode uint32, _ *ep.Poller) {
		//log.Printf("sock active \n")
		//if eventCode&syscall.EPOLLRDHUP != 0 {
		//	log.Printf("client closed\n")
		//	//log.Printf("conn down\n")
		//	//客户端主动断开连接，调用 epoll_remove
		//	_ = ep.EpollRemove(epfd, conn, &poller)
		//	return
		//}
		if eventCode&syscall.EPOLLERR != 0 {
			log.Printf("ep error\n")
			_ = ep.EpollRemove(epfd, conn, &poller)
			return
		}
		//log.Printf("receive code %x\n", eventCode)
		//log.Printf("socket active\n")
		nbuf, err := conn.Read(buf)
		if err != nil {
			_ = ep.EpollRemove(epfd, conn, &poller)
			return
		}
		//只要有数据包，就不可能小于0
		if nbuf <= 0 {
			_ = ep.EpollRemove(epfd, conn, &poller)
			return
		}
		nw, err := handler.mp[conn].Write(buf[:nbuf])
		if err != nil {
			log.Printf("write error %v", err)
			_ = ep.EpollRemove(epfd, conn, &poller)
			return
		}
		if nw <= 0 {
			_ = ep.EpollRemove(epfd, conn, &poller)
			return
		}

	}

	var queue = make([]syscall.EpollEvent, consts.TcpEpollEventQueueSize)

	_, err := ep.Listen(&poller, hostname, iport, consts.TcpBackLog, queue[:])
	if err != nil {
		log.Printf("epoll error %v", err)
	}

}
