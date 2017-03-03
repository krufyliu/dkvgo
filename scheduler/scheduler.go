package scheduler

import (
	"net"
	"log"
	"github.com/krufyliu/dkvgo/task"
	"sync"
	"time"
)

// DkvScheduler d
type DkvScheduler struct {
	sync.WaitGroup
	opts        *Options
	tcpListener net.Listener
	Store 		TaskStore
	Pool        *WorkerPool
	RunningTasks map[int]*task.Task
}

func (s *DkvScheduler) Main() {
	tcpListener, err := net.Listen("tcp", s.opts.TCPAddr)
	if err != nil {
		log.Fatalf("FATAL: listen %s failed - %s", s.opts.TCPAddr， err)
	}
	s.tcpListener = tcpListener
	s.runTcpServer()
	s.schedTasks()
	s.Wait()
}

func (s *DkvScheduler) schedTasks() {
	s.Add(1)
	defer s.Done()
	for {
		if !s.Pool.HasIdleWorker() {
			time.Sleep(2*time.Second)
			continue
		}
		t := s.Store.GetTask()
		if t == nil {
			time.Sleep(5*time.Second)
			continue
		}
	}
}

func (s *DkvScheduler) runTcpServer() {
	s.Add(1)
	defer s.Done()
	TCPServer(s.tcpListener, ProtocolLoop{ctx: s})
}