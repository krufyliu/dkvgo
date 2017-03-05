package scheduler

import (
	"github.com/krufyliu/dkvgo/job"
	"container/list"
	"sync"
)

type TaskPool struct {
	sync.Mutex
	ctx *DkvScheduler
	queue *list.List
}

func (tp *TaskPool) GetTask() *job.Task {
	tp.Lock()
	defer tp.Unlock()
	var elem = tp.queue.Front()
	if elem == nil && !tp.tryFill() {
		return nil
	}
	elem = tp.queue.Front()
	if elem != nil {
		return elem.Value.(*job.Task)
	}
	return nil
}

func (tp *TaskPool) tryFill() bool{
	var _job = tp.ctx.Store.GetJob()
	if _job == nil {
		return false
	}
	tp.ctx.Store.LoadJobState(_job)
	_job.Map()
	for _, opt := range(_job.TaskOpts) {
		if opt.FrameAt != opt.EndFrame {
			tp.queue.PushBack(&job.Task{Job: _job, Options: opt})
		}
	}
	return true
}