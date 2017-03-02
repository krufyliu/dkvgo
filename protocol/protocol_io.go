package protocol

import (
	"net"
)

type ProtocolIO interface{
    IOLoop(conn *net.Conn) error
}