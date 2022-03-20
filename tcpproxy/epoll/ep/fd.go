package ep

import (
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
	return syscall.Write(int(fd), buf)
}

func (c SockFd) String() string {
	return fmt.Sprintf("Connection{fd=%d}", c)
}

func (c SockFd) Read(buf []byte) (int, error) {
	return syscall.Read(int(c), buf)
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
