package scheduler

import (
	"container/list"
	"log"
	"sync"

	"github.com/krufyliu/dkvgo/job"
)

type TaskPool struct {
	sync.Mutex
	ctx   *DkvScheduler
	queue *list.List
}

func NewTaskPool(ctx *DkvScheduler) *TaskPool {
	return &TaskPool{
		ctx:   ctx,
		queue: list.New(),
	}
}

func (tp *TaskPool) GetTask() *job.Task {
	tp.Lock()
	defer tp.Unlock()
	var elem *list.Element
	for elem = tp.queue.Front(); elem != nil; elem = elem.Next() {
		status := elem.Value.(*job.Task).Job.GetStatus()
		if status == 0x01 || status == 0x02 {
			break
		}
	}
	if elem != nil {
		var task = elem.Value.(*job.Task)
		// add job's running task
		task.Job.IncRunning()
		return task
	}
	if !tp.tryFill() {
		return nil
	}
	elem = tp.queue.Front()
	if elem != nil {
		var task = elem.Value.(*job.Task)
		// add job's running task
		task.Job.IncRunning()
		return task
	}
	return nil
}

func (tp *TaskPool) PushFront(task *job.Task) {
	tp.Lock()
	defer tp.Unlock()
	task.Job.DecRunning()
	tp.queue.PushFront(task)
}

func (tp *TaskPool) tryFill() bool {
	var _job = tp.ctx.Store.GetJob()
	if _job == nil {
		log.Printf("no job to fill task pool\n")
		return false
	}
	log.Printf("fill task pool with %d\n", _job.ID)
	tp.ctx.Store.LoadJobState(_job)
	_job.Init()
	// set state to accept
	_job.Status = 0x01
	tp.ctx.Store.UpdateJob(_job)
	for _, opt := range _job.TaskOpts {
		if opt.FrameAt != opt.EndFrame {
			tp.queue.PushBack(&job.Task{Job: _job, Options: opt})
		}
	}
	return true
}
