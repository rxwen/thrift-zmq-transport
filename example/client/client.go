package main

import (
	"fmt"
	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/rxwen/thrift-zmq-transport"
	"github.com/rxwen/thrift-zmq-transport/example/gen-go/demo/rpc"
	//"net"
	//"os"
)

func main() {
	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()

	//transport, err := thrift.NewTSocket(net.JoinHostPort("127.0.0.1", "19092"))
	transport := zmqtransport.NewTDealerTransport("tcp://127.0.0.1:19093")

	useTransport := transportFactory.GetTransport(transport)
	useTransport.Open()
	client := rpc.NewRpcServiceClientFactory(useTransport, protocolFactory)

	defer transport.Close()

	r1, e1 := client.Foo(42, "", "owner")
	fmt.Println("Call Foo: ", r1, e1)
}
