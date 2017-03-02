package scheduler

import (
	"net"
)

type ProtocolLoop struct {
	ctx *DkvScheduler
}

func (loop *DkvScheduler) Handle(conn net.Conn) {
}
