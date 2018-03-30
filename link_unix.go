// +build darwin dragonfly freebsd linux nacl netbsd openbsd solaris

package gotfo

import (
	"net"
	"syscall"
	_ "unsafe"
)

//go:linkname connectFunc net.connectFunc
func connectFunc(int, syscall.Sockaddr) error

//go:linkname testHookCanceledDial net.testHookCanceledDial
func testHookCanceledDial()

//go:linkname getsockoptIntFunc net.getsockoptIntFunc
func getsockoptIntFunc(int, int, int) (int, error)

//go:linkname wrapSyscallError net.wrapSyscallError
func wrapSyscallError(string, error) error

//go:linkname addrFunc net.addrFunc
func (fd *netFD) addrFunc() func(syscall.Sockaddr) net.Addr
