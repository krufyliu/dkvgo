package scheduler

import (
	"net"
	"sync"

	"github.com/krufyliu/dkvgo/task"
)

type Worker struct {
	conn       net.Conn
	remoteAddr string
	lastUpdate int64
	busy       bool
	ts         *task.TaskSegment
}

func (worker *Worker) Attach(ts *task.TaskSegment) {
	worker.ts = ts
}

func (worker *Worker) Dettach() *task.TaskSegment {
	var ts = worker.ts
	worker.ts = nil
	return ts
}

type WorkerPool struct {
	sync.Mutex
	workers map[string]*Worker
}

func (pool *WorkerPool) Add(w *Worker) {
	pool.Lock()
	pool.workers[w.RemoteAddr] = w
	pool.Unlock()
}

func (pool *WorkerPool) Remove(w *Worker) {
	pool.Lock()
	delete(pool.workers, w.remoteAddr)
	pool.Unlock()
}
