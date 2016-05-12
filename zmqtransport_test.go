package zmqtransport_test

import (
	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/rxwen/thrift-zmq-transport"
	"testing"
)

func TestInterfaceImplementation(t *testing.T) {
	var _ thrift.TTransport = zmqtransport.TDealerTransport{}
	var _ thrift.TTransport = (*zmqtransport.TDealerTransport)(nil)
	var _ thrift.TServerTransport = zmqtransport.TRouterServerTransport{}
	var _ thrift.TServerTransport = (*zmqtransport.TRouterServerTransport)(nil)
	var _ thrift.TTransport = zmqtransport.TRouterTransport{}
	var _ thrift.TTransport = (*zmqtransport.TRouterTransport)(nil)
}
