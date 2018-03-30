# TCP Fast Open in Go
TCP Fast Open on Windows 10 (since version 1607) and Linux (since 3.7), go1.8/1.9 only.

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
