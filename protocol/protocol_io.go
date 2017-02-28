package protocol

import (
	"net"
)

type interface ProtocolIO {
    func IOLoop(conn *net.Conn) error
}