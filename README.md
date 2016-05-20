[![Build Status](https://travis-ci.org/rxwen/thrift-zmq-transport.svg?branch=master)](https://travis-ci.org/rxwen/thrift-zmq-transport)


# thrift-zmq-transport
zmq [transport](https://github.com/apache/thrift/blob/master/lib/go/thrift/transport.go) for thrift framework

## supported zmq socket

- DEALER / ROUTER

## cross compile on mac os x 

Becase the [zmq4](https://github.com/pebbe/zmq4) makes use of CGO, it's not possible to cross compile a linux binary on mac os platform. Use [golang-zeromq](https://hub.docker.com/r/rxwen/golang-zeromq/) docker image to cross compile. Usage:
	docker run --rm -w /go/src/path/to/go/project -v ${GOPATH}:/go rxwen/golang-zeromq go build ./

