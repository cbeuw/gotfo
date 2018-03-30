package gotfo

import (
	"net"
	"syscall"
	_ "unsafe"
)

//go:linkname connectFunc net.connectFunc
func connectFunc(syscall.Handle, syscall.Sockaddr) error

//go:linkname sysSocket net.sysSocket
func sysSocket(int, int, int) (syscall.Handle, error)

//go:linkname sysInit net.sysInit
func sysInit()

//go:linkname ok net.(*conn)ok
func (*conn) ok() bool

//go:linkname wrapSyscallError net.wrapSyscallError
func wrapSyscallError(string, error) error

//go:linkname addrFunc net.(*netFd)addrFunc
func (fd *netFD) addrFunc() func(syscall.Sockaddr) net.Addr
