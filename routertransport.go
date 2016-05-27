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

// TRouterServerTransport is a zeromq Router based TServerTransport implementation.
type TRouterServerTransport struct {
	transports map[string]TRouterTransport
	sock       *zmq.Socket
	endpoint   string
	cwrite     chan WriteMessage
	caccept    chan []byte
}

// NewTRouterServerTransport function instantiates a new TRouterServerTransport to specified endpoint
func NewTRouterServerTransport(endpoint string) TRouterServerTransport {
	fmt.Println("new TRouterServerTransport")
	sock, err := zmq.NewSocket(zmq.ROUTER)

	if err != nil {
		panic(err.Error())
	}

	cwrite := make(chan WriteMessage, 100)
	caccept := make(chan []byte, 100)
	go func() {
		for {
			msg := <-cwrite
			sock.SendMessage(msg.id, msg.data)
			//fmt.Println("send ", msg)
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
			msg, err := sock.RecvMessageBytes(0)
			//fmt.Println("number of peers ", len(t.transports))
			id := msg[0]
			var transport TRouterTransport
			if err == nil {
				buf := msg[1]

				if trans, ok := t.transports[string(id)]; !ok {
					cread := make(chan []byte, 100)
					transport = TRouterTransport{
						id:     id,
						cread:  cread,
						cwrite: t.cwrite,
						server: &t,
					}
					fmt.Println("create transport channel for", transport)

					t.transports[string(id)] = transport
					caccept <- id
				} else {
					transport = trans
				}
				transport.cread <- buf
			}
		}
	}()
	return t
}

func (t TRouterServerTransport) Listen() error {
	t.sock.Bind(t.endpoint)
	fmt.Println("TRouterServerTransport.Listen ", t.endpoint)
	return nil
}

func (t TRouterServerTransport) Accept() (thrift.TTransport, error) {
	id := <-t.caccept
	fmt.Println("TRouterServerTransport.Accept ", id, t.transports[string(id)])
	return t.transports[string(id)], nil
}

func (t TRouterServerTransport) Close() error {
	fmt.Println("TRouterServerTransport.Close")

	transports := make([]TRouterTransport, 0)
	for _, trans := range t.transports {
		transports = append(transports, trans)
	}
	for _, trans := range transports {
		trans.Close()
	}
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
	server *TRouterServerTransport
}

func (t TRouterTransport) Open() error {
	return nil
}

func (t TRouterTransport) IsOpen() bool {
	return true
}

func (t TRouterTransport) Close() error {
	fmt.Println("TRouterTransport.Close")
	//t.cwrite <- WriteMessage{t.id, []byte{}} // send 0 bytes indicates close
	delete(t.server.transports, string(t.id))
	return nil
}

func (t TRouterTransport) Read(buf []byte) (l int, err error) {
	//fmt.Println("TRouterTransport.Read wait, len ", len(buf), buf)
	data := <-t.cread
	copy(buf, data)
	return len(data), nil
}

func (t TRouterTransport) Write(buf []byte) (int, error) {
	data := make([]byte, len(buf))
	copy(data, buf)
	msg := WriteMessage{
		id:   t.id,
		data: data,
	}

	t.cwrite <- msg
	//fmt.Println("write ", msg)
	return len(buf), nil
}

func (p TRouterTransport) Flush() error {
	return nil
}

func (p TRouterTransport) RemainingBytes() (num_bytes uint64) {
	//fmt.Println("RemainingBytes, hard code to 128")
	return uint64(4096)
}
