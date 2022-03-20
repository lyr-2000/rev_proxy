package ep

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"syscall"
)

//const (
//	readEventCode  uint32 = 1 << 0
//	writeEventCode        = 1 << 1
//	errEventCode          = 1 << 2
//)

type Poller struct {
	SockFd       SockFd
	EpFd         EpollFd
	OnFdRemoved  func(epfd EpollFd, conn SockFd)
	OnConnOpen   func(epfd EpollFd, conn SockFd)
	OnMsgReceive func(epfd EpollFd, conn SockFd, eventCode uint32)
}

func (p *Poller) Close() error {
	p.EpFd.Close()
	p.SockFd.Close()
	return nil
}

func EpollRemove(fd EpollFd, sockFd SockFd, poller *Poller) error {
	defer func() {
		if poller != nil {
			poller.OnFdRemoved(fd, sockFd)
		}
	}()
	if err := syscall.EpollCtl(int(fd), syscall.EPOLL_CTL_DEL, int(sockFd), nil); err != nil {
		fmt.Printf("epoll_ctl remove err: %+v\n", err)
		return err
	}
	return nil
}
func EpollCtl(epfd EpollFd, fd SockFd, eventsListen uint32) error {
	var event syscall.EpollEvent
	//监听 sock 连接事件
	event.Fd = int32(fd)
	//epollin 事件
	event.Events = eventsListen
	if e := syscall.EpollCtl(int(epfd), syscall.EPOLL_CTL_ADD, int(fd), &event); e != nil {
		log.Printf("epoll_ctrl error %v\n", e)
		return e
	}
	return nil
}

func Listen(poll *Poller, host string, port int, backlog int, eventQueue []syscall.EpollEvent) (*Poller, error) {
	//if poll == nil {
	//	panic("empty poller")
	//}
	if poll == nil {
		poll = &Poller{}
	}
	if host == "" {
		host = "0.0.0.0"
	}
	if backlog <= 0 {
		//backlog参数
		backlog = syscall.SOMAXCONN
	}
	fd, err := syscall.Socket(syscall.AF_INET, syscall.O_NONBLOCK|syscall.SOCK_STREAM, 0) //创建socket
	if err != nil || fd < 0 {
		log.Printf("%+v\n", err)
		os.Exit(1)
	}
	defer poll.Close()
	poll.SockFd = SockFd(fd)
	//set nonblock mode
	poll.SockFd.SetNonBlock()

	addr := syscall.SockaddrInet4{Port: port}
	copy(addr.Addr[:], net.ParseIP(host).To4())
	syscall.Bind(fd, &addr)     //绑定ip,端口。
	syscall.Listen(fd, backlog) //监听端口fd，  n表示 backlog 参数，我设置为最大

	epfd, err := syscall.EpollCreate1(0)
	if err != nil || epfd < 0 {
		poll.Close()
		return nil, errors.New("epoll create error")
	}
	poll.EpFd = EpollFd(epfd)
	//epollfd 监听 sockfd 连接事件
	err = EpollCtl(poll.EpFd, poll.SockFd, syscall.EPOLLIN)
	if err != nil {
		poll.Close()
		return nil, err
	}
	//var acceptEvent syscall.EpollEvent
	for {
		nevent, err := syscall.EpollWait(epfd, eventQueue[:], -1)
		//log.Printf("nevent %v, %v", nevent, err)
		if err != nil {
			if err == syscall.EINTR {
				continue
			}
			fmt.Printf("epoll_wait: %+v\n", err)
			//break
			continue
		}
		for ev := 0; ev < nevent; ev++ {
			if int(eventQueue[ev].Fd) == fd {
				connFd, _, err := syscall.Accept(fd)
				if err != nil {
					log.Printf("accept error %+v\n", err)
					continue
				}
				sockFd := SockFd(connFd)
				sockFd.SetNonBlock()
				//监听 连接本断开的 sock 的 各种发送消息的事件
				if err := EpollCtl(poll.EpFd, sockFd, syscall.EPOLLIN|syscall.EPOLLPRI); err != nil {
					log.Printf("error at epoll ctrl %+v\n", err)
				} else {
					//无异常
					poll.OnConnOpen(poll.EpFd, sockFd)
				}

			} else {

				poll.OnMsgReceive(poll.EpFd, SockFd(eventQueue[ev].Fd), eventQueue[ev].Events)

			}
		}

	}

	return poll, nil

}
