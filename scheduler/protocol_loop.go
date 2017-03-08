package scheduler

import (
	"net"
)

type ProtocolLoop struct {
	ctx        *DkvScheduler
}

func (loop *ProtocolLoop) Handle(conn net.Conn) {
	var worker = Worker{
		ctx: loop.ctx, 
	}
	worker.IOLoop(conn)
}




