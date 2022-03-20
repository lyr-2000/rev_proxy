package epoll

//
//const (
//	//EPOLLET        = syscall.EPOLLPRI | syscall.EPOLLIN //这里只是监听读事件
//	MaxEpollEvents = 32
//)
//
//var buf [8]byte
//
//func echo(epollfd, fd int) {
//	//defer syscall.Close(fd)
//	//for {
//	nbytes, e := syscall.Read(fd, buf[:])
//	if nbytes > 0 {
//		fmt.Printf(">>> %s", buf)
//		_, e := syscall.Write(fd, buf[:nbytes])
//		if e != nil {
//			log.Printf("error write [%+v]", e)
//			//panic(e)
//		}
//		fmt.Printf("<<< %s", buf)
//	}
//	if e != nil {
//		log.Printf("error %v\n", e)
//	}
//	if nbytes == 0 {
//		//如果无法读取，说明客户端被关闭了，要移除连接
//		if err := syscall.EpollCtl(epollfd, syscall.EPOLL_CTL_DEL, fd, nil); err != nil {
//			fmt.Println("epoll_ctl: ", err)
//			//os.Exit(1)
//		}
//		log.Printf("close epoll fd %v\n", fd)
//
//	}
//	//}
//}
//func main() {
//	var event syscall.EpollEvent
//	var events [MaxEpollEvents]syscall.EpollEvent
//	fd, err := syscall.Socket(syscall.AF_INET, syscall.O_NONBLOCK|syscall.SOCK_STREAM, 0) //创建socket
//	if err != nil {
//		fmt.Println(err)
//		os.Exit(1)
//	}
//	defer syscall.Close(fd)
//	if err = syscall.SetNonblock(fd, true); err != nil { //设置非阻塞模式
//		fmt.Println("setnonblock1: ", err)
//		os.Exit(1)
//	}
//	addr := syscall.SockaddrInet4{Port: 50000}
//	copy(addr.Addr[:], net.ParseIP("0.0.0.0").To4())
//	syscall.Bind(fd, &addr) //绑定ip,端口。
//	syscall.Listen(fd, 10)  //监听端口
//	epfd, e := syscall.EpollCreate1(0)
//	if e != nil {
//		fmt.Println("epoll_create1: ", e)
//		os.Exit(1)
//	}
//	defer syscall.Close(epfd)
//	event.Events = syscall.EPOLLIN
//	event.Fd = int32(fd) //设置监听描述符
//	if e = syscall.EpollCtl(epfd, syscall.EPOLL_CTL_ADD, fd, &event); e != nil {
//		fmt.Println("epoll_ctl: ", e)
//		os.Exit(1)
//	}
//	log.Printf("ok=>enter\n")
//	for {
//		nevents, e := syscall.EpollWait(epfd, events[:], -1) //获取就绪事件
//		if e != nil {
//			if e == syscall.EINTR {
//				continue
//			}
//			fmt.Printf("epoll_wait: %+v\n", e)
//			//break
//			continue
//		}
//		fmt.Printf("nevent cnt = %d\n", nevents)
//		for ev := 0; ev < nevents; ev++ {
//			if int(events[ev].Fd) == fd {
//				connFd, _, err := syscall.Accept(fd) //接受请求
//				if err != nil {
//					fmt.Println("accept: ", err)
//					continue
//				}
//				syscall.SetNonblock(fd, true)
//				//syscall.EPOLLIN
//				//监听读事件, 模式是水平触发
//				event.Events = syscall.EPOLLIN | syscall.EPOLLPRI
//				event.Fd = int32(connFd)
//				if err := syscall.EpollCtl(epfd, syscall.EPOLL_CTL_ADD, connFd, &event); err != nil {
//					fmt.Print("epoll_ctl: ", connFd, err)
//					os.Exit(1)
//				}
//			} else {
//				//var rEvents uint32
//				if ((events[ev].Events & unix.POLLHUP) != 0) && ((events[ev].Events & unix.POLLIN) == 0) {
//					//rEvents |= EventErr
//					log.Printf("event err\n")
//				}
//				if (events[ev].Events&unix.EPOLLERR != 0) || (events[ev].Events&unix.EPOLLOUT != 0) {
//					log.Printf("write event\n")
//				}
//				if events[ev].Events&(unix.EPOLLIN|unix.EPOLLPRI|unix.EPOLLRDHUP) != 0 {
//					log.Printf("read event")
//				}
//				log.Printf("events.fd=%d\n", events[ev].Fd)
//				log.Printf("events %#v", events[ev])
//				/*go*/ echo(epfd, int(events[ev].Fd))
//			}
//		}
//	}
//}
