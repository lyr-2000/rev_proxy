//go:build linux
// +build linux

package tcpproxy

func Default() TcpReverseProxy {
	//if runtime.GOOS == "windows" {
	//	//可以利用 iocp优化
	//}

	if runtime.GOOS == "linux" {
		//log.Printf("use epoll handler\n")
		//使用 epoll
		//return &v1.EpollHandler{}
		log.Printf("use epoll features \n")
		return &epoll.EpHandler{}
	}
	//通用的方法处理连接
	return &commonhandle.CommonHandler{}
}
