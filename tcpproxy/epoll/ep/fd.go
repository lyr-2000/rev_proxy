package ep

import (
	"errors"
	"fmt"
	"net"
	"reflect"
	"syscall"
)

type SockFd int
type EpollFd int

func (fd *SockFd) SetNonBlock() {
	syscall.SetNonblock(int(*fd), true)
}

func (fd *SockFd) Close() error {
	if *fd == 0 {
		return nil
	}
	mfd := int(*fd)
	e := syscall.Close(mfd)

	//*fd = 0
	return e
}

func (fd *EpollFd) Close() error {
	if *fd == 0 {
		return nil
	}
	mfd := int(*fd)
	e := syscall.Close(mfd)

	//*fd = 0
	return e
}

func (fd SockFd) Write(buf []byte) (int, error) {
	//if fd == 0 {
	//	return 0,nil
	//}
	return syscall.Write(int(fd), buf)
}

func (c SockFd) String() string {
	return fmt.Sprintf("Connection{fd=%d}", c)
}

func (c SockFd) Read(buf []byte) (int, error) {
	return syscall.Read(int(c), buf)
}

func GetSockFD(conn net.Conn) int {
	return sockFD(conn)
}
func sockFD(conn net.Conn) int {
	tcpConn := reflect.Indirect(reflect.ValueOf(conn)).FieldByName("conn")
	//if tls {
	//	tcpConn = reflect.Indirect(tcpConn.Elem())
	//}
	fdVal := tcpConn.FieldByName("fd")
	pfdVal := reflect.Indirect(fdVal).FieldByName("pfd")

	return int(pfdVal.FieldByName("Sysfd").Int())
}
func Open(network string, addr string) (SockFd, net.Conn, error) {
	if network == "" {
		network = "tcp"
	}
	dial, err := net.Dial(network, addr)
	if err != nil {
		return 0, dial, err
	}
	fd := sockFD(dial)
	return SockFd(fd), dial, err

}

func GetSockname(fd SockFd) (syscall.Sockaddr, error) {
	lsa, err := syscall.Getpeername(int(fd))
	//net.ParseIP("0.0.0.0").To4()
	//
	//lsa.(syscall.SockaddrInet4)
	return lsa, err

}
func GetIpv4(sockaddr syscall.Sockaddr) (string, error) {
	switch sockaddr.(type) {
	case *syscall.SockaddrInet4:
		s := sockaddr.(*syscall.SockaddrInet4)
		return fmt.Sprintf("%d.%d.%d.%d", s.Addr[0], s.Addr[1], s.Addr[2], s.Addr[3]), nil
	default:

		return "", errors.New("not ipv4 address")
	}

}
