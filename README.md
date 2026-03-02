## gomulticast

This package provides a simple wrapper over the UDP multicast functions in [`net`](https://pkg.go.dev/net). One of the primary benefits of using this package is the ability to mock network access, allowing packages to easily test the implementation of their network code as-is:

```golang
import "github.com/nitroshare/gomulticast"

// Get a slice of network interfaces available
ifaces, err := gomulticast.Interfaces()
if err != nil {
    panic(err)
}

// Find the first one that can multicast
var iface gomulticast.Interface
for _, i := range ifaces {
    if i.Interface().Flags & net.FlagMulticast != 0 {
        iface = i
        break
    }
}
if iface == nil {
    panic("no multicast interfaces found")
}

// Multicast address
addr := &net.UDPAddr{
    IP:   net.IP([]byte{1, 2, 3, 4}),
	Port: 1234,
}

// Create a listener for sending and receiving packets
l, err := gomulticast.New("udp4", iface, addr)
if err != nil {
    panic(err)
}
defer l.Close()

// Send a packet
l.Write(&gomulticast.Packet{
    Addr: addr,
    Data: []byte("test data"),
})

// Read a packet (this blocks)
p, _ := l.Read()
```

### Mock Interfaces

If for example you wanted to test the code above without modifying it or relying on the host's network stack, gomulticast has you covered (no pun intended)!

```golang
// Mock everything in the package
gomulticast.Mock()
defer gomulticast.Unmock()

// Create a mocked interface and make future calls to Interfaces() return it
i := gomulticast.NewMockInterface()
gomulticast.AddMockInterface(i)

// Now whenever Interfaces() is called, this will be the only one returned

// Packets can be queued for reading from it...
i.QueueForRead(&gomulticast.Packet{...})

// ...and packets written to it can be dequeued
p, _ := i.DequeueWrite()
```
