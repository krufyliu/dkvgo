package scheduler

import (
	"net"
	"log"
)

// DkvScheduler d
type DkvScheduler struct {
	opts        *Options
	tcpListener net.Listener
	Pool        *WorkerPool
}

func (s *DkvScheduler) Main() {
	tcpListener, err := net.Listen("tcp", s.opts.TCPAddr)
	if err != nil {
		log.Fatalf("FATAL: listen %s failed - %s", s.opts.TCPAddr， err)
	}
	s.tcpListener = tcpListener
}