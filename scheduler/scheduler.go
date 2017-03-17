package scheduler

import (
	"log"
	"net"
	"sync"

	"github.com/krufyliu/dkvgo/job"
	"github.com/krufyliu/dkvgo/job/store"
	"github.com/krufyliu/dkvgo/scheduler/tracker"
)

// DkvScheduler d
type DkvScheduler struct {
	sync.WaitGroup
	mu          sync.Mutex
	opts        *_options
	tcpListener net.Listener
	TaskPool    *TaskPool
	Store       store.JobStore
	runningJobs map[int]*job.Job
}

func newDkvScheduler() *DkvScheduler {
	var sched = &DkvScheduler{
		opts: Options,
		//Store: store.NewMockStore(),
		Store: store.NewDatabaseStore(Options.DBType, Options.DBAddr),
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
	tracker.InitWithStore(s.Store)
	tcpListener, err := net.Listen("tcp", s.opts.TCPAddr)
	if err != nil {
		log.Fatalf("FATAL: listen %s failed - %s\n", s.opts.TCPAddr, err)
	}
	s.tcpListener = tcpListener
	log.Printf("TCP listen on %s\n", tcpListener.Addr())
	s.Add(1)
	s.Add(1)
	go s.runTcpServer()
	go s.runApiServer()
	s.Wait()
}

func (s *DkvScheduler) runTcpServer() {
	defer s.Done()
	TCPServer(s.tcpListener, &ProtocolLoop{ctx: s})
}

func (s *DkvScheduler) runApiServer() {
	defer s.Done()
	log.Printf("HTTP listen on %s", s.opts.HTTPAddr)
	APIServer(s.opts.HTTPAddr).ListenAndServe()
}
