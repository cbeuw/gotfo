# Usage
```go
// dial an address with data, it returns a net.Conn
conn, err := gotfo.Dial(address, true, data)
// or dial without fast open
conn, err := gotfo.Dial(address, false, nil)

// listen with fast open, it returns a net.Listener
listener, err := gotfo.Listen(address, true)
// or listen without fast open
listener, err := gotfo.Listen(address, false)
```

# Change notes:

./fd_unix.go and ./fd_windows.go are copied from package "net"

./internal/* are from internal

## From "net":

fd_unix.go: import

fd_windows.go: function `connect`. Added `data []byte` as an extra argument and uses ConnectEx to dial and send data.

## From "internal/poll":

fd_windows.go: function `ConnectEx`. Added `data []byte` as an extra argument and dial and sends data with `ConnectExFunc`.

fd_poll_runtime.go: remove abstract function declarations for go:linkname

fd_mutex.go: remove abstract function declarations for go:linkname

## From "internal/syscall/windows":

zsyscall_windows.go: import

