// +build dragonfly freebsd linux nacl netbsd openbsd solaris

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

func socket(family int, fastOpen bool) (int, error) {
	fd, err := syscall.Socket(family, syscall.SOCK_STREAM, 0)
	if err != nil {
		return 0, err
	}
	if fastOpen {
		if err := syscall.SetsockoptInt(fd, syscall.IPPROTO_TCP, TCP_FASTOPEN, 1); err != nil {
			return 0, err
		}
	}

	if err := syscall.SetsockoptInt(fd, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1); err != nil {
		return 0, err
	}

	if err := syscall.SetsockoptInt(fd, syscall.IPPROTO_TCP, syscall.TCP_NODELAY, 1); err != nil {
		return 0, err
	}

	return fd, nil
}

func Listen(address string, fastOpen bool) (net.Listener, error) {
	laddr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return nil, err
	}

	fd, err := socket(syscall.AF_INET, fastOpen)
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

	nfd, _ := newFD(fd, syscall.AF_INET, syscall.SOCK_STREAM, "tcp")
	if err := nfd.init(); err != nil {
		syscall.Close(fd)
		return nil, err
	}

	return newTCPListener(nfd, false), nil
}

func Dial(address string, fastOpen bool, data []byte) (*net.TCPConn, error) {
	return DialContext(context.Background(), address, fastOpen, data)
}

var fdCallback func(int)

func SetFdCallback(fn func(int)) {
	fdCallback = fn
}

func DialContext(ctx context.Context, address string, fastOpen bool, data []byte) (*net.TCPConn, error) {
	raddr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return nil, err
	}

	fd, err := socket(syscall.AF_INET, fastOpen)
	if err != nil {
		syscall.Close(fd)
		return nil, err
	}

	sa := tcpAddrToSockaddr(raddr)

	nfd, _ := newFD(fd, syscall.AF_INET, syscall.SOCK_STREAM, "tcp")
	if err := nfd.init(); err != nil {
		syscall.Close(fd)
		return nil, err
	}

	if fdCallback != nil {
		fdCallback(nfd.pfd.Sysfd)
	}

	for {
		if fastOpen {
			// MSG_FASTOPEN=0x20000000
			// syscall.MSG_FASTOPEN is not defined on mac
			err = syscall.Sendto(nfd.pfd.Sysfd, data, 0x20000000, sa)
		} else {
			err = syscall.Connect(nfd.pfd.Sysfd, sa)
		}
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
