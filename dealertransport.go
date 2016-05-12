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

func NewTDealerTransport(endpoint string) TDealerTransport {
	sock, _ := zmq.NewSocket(zmq.DEALER)

	return TDealerTransport{
		sock:     sock,
		endpoint: endpoint,
	}
}

func (t TDealerTransport) Open() error {
	fmt.Println("Open", t.endpoint)
	return t.sock.Connect(t.endpoint)
}

func (t TDealerTransport) IsOpen() bool {
	fmt.Println("IsOpen")
	return true
}

func (t TDealerTransport) Close() error {
	fmt.Println("Close")
	c := make([]byte, 0)
	t.sock.SendBytes(c, 0) // send 0 bytes indicates close
	return t.sock.Close()
}

func (t TDealerTransport) Read(buf []byte) (l int, err error) {
	//fmt.Println("about to Read num of bytes ", len(buf))
	data, err := t.sock.RecvBytes(0)
	//fmt.Println("received num of bytes ", len(data), data)
	l = len(data)
	copy(buf, data)
	return l, err
}

func (t TDealerTransport) Write(buf []byte) (int, error) {
	//fmt.Println("Write ", buf)
	return t.sock.SendBytes(buf, 0)
}

func (p TDealerTransport) Flush() error {
	//fmt.Println("Flush")
	return nil
}

func (p TDealerTransport) RemainingBytes() (num_bytes uint64) {
	//fmt.Println("RemainingBytes, hard code to 128")
	return uint64(4096)
}
