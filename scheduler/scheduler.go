package scheduler

import (
	"net"
	"log"
	"sync"
	"github.com/krufyliu/dkvgo/job"
)

// DkvScheduler d
type DkvScheduler struct {
	sync.WaitGroup
	opts        *Options
	tcpListener net.Listener
	httpListener net.Listener
	Store 		TaskStore
	RunningTasks map[int]*job.Job
}

func (s *DkvScheduler) Main() {
	tcpListener, err := net.Listen("tcp", s.opts.TCPAddr)
	if err != nil {
		log.Fatalf("FATAL: listen %s failed - %s", s.opts.TCPAddr, err)
	}
	s.tcpListener = tcpListener
	s.runTcpServer()
	s.Wait()
}

func (s *DkvScheduler) runTcpServer() {
	s.Add(1)
	defer s.Done()
	TCPServer(s.tcpListener, &ProtocolLoop{ctx: s})
}