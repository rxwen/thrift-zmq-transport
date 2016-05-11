package zmqtransport

import (
	"fmt"

	zmq "github.com/pebbe/zmq4"

	"git.apache.org/thrift.git/lib/go/thrift"
)

type WriteMessage struct {
	id   []byte
	data []byte
}

// TRouterServerTransport is a TServerTransport implementation.
type TRouterServerTransport struct {
	transports map[string]TRouterTransport
	sock       *zmq.Socket
	endpoint   string
	cwrite     chan WriteMessage
	caccept    chan []byte
}

func NewTRouterServerTransport(endpoint string) TRouterServerTransport {
	fmt.Println("new TRouterServerTransport")
	sock, err := zmq.NewSocket(zmq.ROUTER)

	if err != nil {
	}

	cwrite := make(chan WriteMessage, 100)
	caccept := make(chan []byte, 100)
	go func() {
		for {
			msg := <-cwrite
			sock.SendMessage(msg.id, msg.data)
			fmt.Println("send ", msg)
		}
	}()

	t := TRouterServerTransport{
		transports: make(map[string]TRouterTransport),
		sock:       sock,
		endpoint:   endpoint,
		caccept:    caccept,
		cwrite:     cwrite,
	}
	go func() {
		for {
			id, err := sock.RecvBytes(0)
			var transport TRouterTransport
			if err == nil {
				buf, _ := sock.RecvBytes(0)

				if trans, ok := t.transports[string(id)]; !ok {
					c := make(chan []byte, 1000)
					transport = TRouterTransport{
						id:     id,
						cread:  c,
						cwrite: t.cwrite,
					}
					fmt.Println("create transport channel for", transport)

					t.transports[string(id)] = transport
					//fmt.Println("signal accept ", id, caccept)
					caccept <- id
					fmt.Println("signal accept ok", id)
				} else {
					transport = trans
				}
				//fmt.Println("put data to transport channel ", buf, transport)
				transport.cread <- buf
				//fmt.Println("put data to transport channel ok", buf, transport)
			}
		}
	}()
	return t
}

func (t TRouterServerTransport) Listen() error {
	t.sock.Bind(t.endpoint)
	fmt.Println("TRouterServerTransport.Listen")
	return nil
}

func (t TRouterServerTransport) Accept() (thrift.TTransport, error) {
	id := <-t.caccept
	fmt.Println("TRouterServerTransport.Accept ", id, t.transports[string(id)])
	return t.transports[string(id)], nil
}

func (t TRouterServerTransport) Close() error {
	fmt.Println("TRouterServerTransport.Close")
	return nil
}

func (t TRouterServerTransport) Interrupt() error {
	fmt.Println("TRouterServerTransport.Interrupt")
	return nil
}

// TRouterTransport is a TTransport implementation.
type TRouterTransport struct {
	id     []byte
	cread  chan []byte
	cwrite chan WriteMessage
}

func (t TRouterTransport) GetMapKey() string {
	return string(t.id)
}

//func NewTRouterTransport(c chan []byte) TRouterTransport {
//return TRouterTransport{c: c}
//}

func (t TRouterTransport) Open() error {
	return nil
}

func (t TRouterTransport) IsOpen() bool {
	return true
}

func (t TRouterTransport) Close() error {
	return nil
}

func (t TRouterTransport) Read(buf []byte) (l int, err error) {
	fmt.Println("TRouterTransport.Read wait, len ", len(buf), buf)
	data := <-t.cread
	//fmt.Println("TRouterTransport.Read got, len ", len(data), data)
	for i, b := range data {
		buf[i] = b
	}
	fmt.Println("TRouterTransport.Read ", buf)
	return len(data), nil
	//buf = <-t.cread
	//return len(buf), nil
}

func (t TRouterTransport) Write(buf []byte) (int, error) {
	data := make([]byte, len(buf))
	for i, b := range buf {
		data[i] = b
	}
	msg := WriteMessage{
		id:   t.id,
		data: data,
	}

	t.cwrite <- msg
	fmt.Println("write ", msg)
	return len(buf), nil
}

func (p TRouterTransport) Flush() error {
	return nil
}

func (p TRouterTransport) RemainingBytes() (num_bytes uint64) {
	fmt.Println("RemainingBytes, hard code to 128")
	return uint64(128)
}
