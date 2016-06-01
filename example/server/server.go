package main

import (
	"fmt"
	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/rxwen/thrift-zmq-transport"
	"github.com/rxwen/thrift-zmq-transport/example/gen-go/demo/rpc"
)

type RpcServiceImpl struct {
}

func (this *RpcServiceImpl) Foo(index int64, code string, owner string) (string, error) {
	fmt.Println("-->FunCall:", index, code, owner)
	return code + owner, nil
}

func main() {
	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()

	serverTransport := zmqtransport.NewTRouterServerTransport("tcp://127.0.0.1:19093")
	//serverTransport, err := thrift.NewTServerSocket("127.0.0.1:19093")

	handler := &RpcServiceImpl{}
	processor := rpc.NewRpcServiceProcessor(handler)

	server := thrift.NewTSimpleServer4(processor, serverTransport, transportFactory, protocolFactory)
	server.Serve()
}
