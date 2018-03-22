// +build darwin dragonfly freebsd linux nacl netbsd openbsd solaris

package gotfo

import (
	"context"
	"net"
	"syscall"
)
import (
	"os"
)

const (
	TCP_FASTOPEN   = 23
	LISTEN_BACKLOG = 23
)

type TFOListener struct {
	*net.TCPListener
	fd *netFD
}

func socket(family int) (int, error) {
	fd, err := syscall.Socket(family, syscall.SOCK_STREAM, 0)
	if err != nil {
		return 0, err
	}

	if err := syscall.SetsockoptInt(fd, syscall.SOL_TCP, TCP_FASTOPEN, 1); err != nil {
		return 0, err
	}

	if err := syscall.SetsockoptInt(fd, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1); err != nil {
		return 0, err
	}

	return fd, nil
}

func Listen(address string, fastOpen bool) (net.Listener, error) {
	if !fastOpen {
		return net.Listen("tcp", address)
	}
	laddr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return nil, err
	}

	fd, err := socket(syscall.AF_INET)
	if err != nil {
		syscall.Close(fd)
		return nil, err
	}

	sa := tcpAddrToSockaddr(laddr)

	if err := syscall.Bind(fd, sa); err != nil {
		syscall.Close(fd)
		return nil, err
	}

	if err := syscall.Listen(fd, LISTEN_BACKLOG); err != nil {
		syscall.Close(fd)
		return nil, err
	}

	nfd := newFD(fd)
	if err := nfd.init(); err != nil {
		syscall.Close(fd)
		return nil, err
	}

	return newTCPListener(nfd, false), nil
}

func Dial(address string, fastOpen bool, data []byte) (*net.TCPConn, error) {
	if fastOpen {
		return DialContext(context.Background(), address, data)
	} else {
		raddr, err := net.ResolveTCPAddr("tcp", address)
		if err != nil {
			return nil, err
		}
		return net.DialTCP("tcp", nil, raddr)
	}
}

var fdCallback func(int, int)

func SetFdCallback(fn func(int, int)) {
	fdCallback = fn
}

func DialContext(ctx context.Context, address string, data []byte) (*net.TCPConn, error) {
	raddr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return nil, err
	}

	fd, err := socket(syscall.AF_INET)
	if err != nil {
		syscall.Close(fd)
		return nil, err
	}

	sa := tcpAddrToSockaddr(raddr)

	nfd := newFD(fd)
	if err := nfd.init(); err != nil {
		syscall.Close(fd)
		return nil, err
	}

	if fdCallback != nil {
		fdCallback(nfd.sysfd, syscall.SOCK_STREAM) // It's always TCP
	}

	for {
		err = syscall.Sendto(nfd.sysfd, data, syscall.MSG_FASTOPEN, sa)
		if err == syscall.EAGAIN {
			continue
		}
		break
	}

	if _, ok := err.(syscall.Errno); ok {
		err = os.NewSyscallError("sendto", err)
	}

	return newTCPConn(nfd), err
}
