package scheduler

import (
	"container/list"
	"net"
	"sync"

	"github.com/krufyliu/dkvgo/job"
)

type Worker struct {
	conn       net.Conn
	remoteAddr string
	lastUpdate int64
	relTask    *job.Task
}

func (worker *Worker) Attach(ts *job.Task) {
	worker.relTask = ts
}

func (worker *Worker) readMonitor() {

}

func (worker *Worker) Dettach() *job.Task {
	var ts = worker.relTask
	worker.relTask = nil
	return ts
}

type WorkerPool struct {
	sync.Mutex
	workers     map[string]*Worker
	idleWorkers *list.List
}

func NewWorkerPool() *WorkerPool {
	return &WorkerPool {
		workers: make(map[string]*Worker),
		idleWorkers: list.New(),
	}
}

func (pool *WorkerPool) Add(w *Worker) {
	pool.Lock()
	pool.workers[w.remoteAddr] = w
	pool.idleWorkers.PushBack(w)
	pool.Unlock()
}

func (pool *WorkerPool) Remove(w *Worker) {
	pool.Lock()
	delete(pool.workers, w.remoteAddr)
	for e := pool.idleWorkers.Front(); e != nil ; e = e.Next() {
		if e.Value.(*Worker) == w {
			pool.idleWorkers.Remove(e)
			break
		}
	}
	pool.Unlock()
}

func (pool *WorkerPool) GetFreeWorker() *Worker {
	pool.Lock()
	defer pool.Unlock()
	var elem = pool.idleWorkers.Front()
	if elem == nil {
		return nil
	}
	return elem.Value.(*Worker)
}

func (pool *WorkerPool) FreeWorker(worker *Worker) {
	pool.Lock()
	pool.idleWorkers.PushBack(worker)
	pool.Unlock()
}

func (pool *WorkerPool) HasIdleWorker() bool {
	pool.Lock()
	defer pool.Unlock()
	return pool.idleWorkers.Len() > 0
}
