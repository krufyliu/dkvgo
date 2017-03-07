package scheduler

import (
	"log"
	"net"
	"sync"

	"github.com/krufyliu/dkvgo/job"
	"github.com/krufyliu/dkvgo/scheduler/store"
)

// DkvScheduler d
type DkvScheduler struct {
	sync.WaitGroup
	mu           sync.Mutex
	opts         *Options
	tcpListener  net.Listener
	httpListener net.Listener
	TaskPool     *TaskPool
	Store        store.JobStore
	runningJobs  map[int]*job.Job
}

func NewDkvScheduler(opts *Options) *DkvScheduler {
	var sched = &DkvScheduler{
		opts:  opts,
		Store: store.NewMockStore(),
	}
	sched.TaskPool = NewTaskPool(sched)
	return sched
}

func (s *DkvScheduler) AddRunningJob(_job *job.Job) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.runningJobs[_job.ID] = _job
}

func (s *DkvScheduler) RemoveRunningJob(_job *job.Job) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.runningJobs, _job.ID)
}

func (s *DkvScheduler) Main() {
	tcpListener, err := net.Listen("tcp", s.opts.TCPAddr)
	if err != nil {
		log.Fatalf("FATAL: listen %s failed - %s\n", s.opts.TCPAddr, err)
	}
	s.tcpListener = tcpListener
	log.Printf("TCP listen on %s\n", tcpListener.Addr().String())
	s.runTcpServer()
	s.Wait()
}

func (s *DkvScheduler) runTcpServer() {
	s.Add(1)
	defer s.Done()
	TCPServer(s.tcpListener, &ProtocolLoop{ctx: s})
}
