package zmqtransport

import (
	"fmt"

	zmq "github.com/pebbe/zmq4"
)

// TDeleaerTransport is a TTransport implementation.
type TDealerTransport struct {
	sock     *zmq.Socket
	endpoint string
}

// NewTDealerTransport function creates a new TDealerTransport instance.
func NewTDealerTransport(endpoint string) TDealerTransport {
	sock, _ := zmq.NewSocket(zmq.DEALER)

	return TDealerTransport{
		sock:     sock,
		endpoint: endpoint,
	}
}

// Open method connects to remote endpoint.
func (t TDealerTransport) Open() error {
	fmt.Println("Open", t.endpoint)
	return t.sock.Connect(t.endpoint)
}

// IsOpen method alwasy return true.
func (t TDealerTransport) IsOpen() bool {
	fmt.Println("IsOpen")
	return true
}

// Close method shuts down the zmq socket.
func (t TDealerTransport) Close() error {
	fmt.Println("Close")
	t.sock.SendBytes([]byte{}, 0) // send 0 bytes indicates close
	return t.sock.Close()
}

// Read method get bytes from socket.
func (t TDealerTransport) Read(buf []byte) (l int, err error) {
	//fmt.Println("about to Read num of bytes ", len(buf))
	data, err := t.sock.RecvBytes(0)
	//fmt.Println("received num of bytes ", len(data), data)
	l = len(data)
	copy(buf, data)
	return l, err
}

// Read method output bytes to socket.
func (t TDealerTransport) Write(buf []byte) (int, error) {
	//fmt.Println("Write ", buf)
	return t.sock.SendBytes(buf, 0)
}

// Flush method performs noop.
func (p TDealerTransport) Flush() error {
	//fmt.Println("Flush")
	return nil
}

// RemainingBytes method returns bytes available.
func (p TDealerTransport) RemainingBytes() (num_bytes uint64) {
	//fmt.Println("RemainingBytes, hard code to 128")
	return uint64(4096)
}
