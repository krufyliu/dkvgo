package scheduler

import (
	"container/list"
	"log"
	"sync"
	"time"

	"github.com/krufyliu/dkvgo/job"
)

type TaskPool struct {
	sync.Mutex
	ctx           *DkvScheduler
	queue         *list.List
	lastQueryTime int64
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
	for elem = tp.queue.Front(); elem != nil; elem = tp.queue.Front() {
		_job := elem.Value.(*job.Task).Job
		status := _job.GetStatus()
		// 找到可以运行的task
		if status == 0x01 || status == 0x02 {
			break
		} else {
			// 如果job处于准备终止状态且当前没有被调度就直接将job置为已停止
			if status == 0x03 && _job.GetRunning() == 0 {
				_job.Status = 0x04
				tp.ctx.Store.UpdateJob(_job)
			}
			tp.queue.Remove(elem)
		}
	}
	if elem != nil {
		var task = elem.Value.(*job.Task)
		// add job's running task
		task.Job.IncRunning()
		tp.queue.Remove(elem)
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
		tp.queue.Remove(elem)
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
	// control the rate of fetch job
	if time.Now().Unix()-tp.lastQueryTime < 10 {
		return false
	}
	var _job = tp.ctx.Store.GetJob()
	if _job == nil {
		log.Printf("no job to fill task pool\n")
		tp.lastQueryTime = time.Now().Unix()
		return false
	}
	log.Printf("fill task pool with %d\n", _job.ID)
	tp.ctx.Store.LoadJobState(_job)
	_job.Init()
	// set state to accept
	_job.Status = 0x01
	tp.ctx.Store.UpdateJob(_job)
	for _, opt := range _job.TaskOpts {
		if opt.FrameAt != opt.EndFrame+1 {
			tp.queue.PushBack(&job.Task{Job: _job, Options: opt})
		}
	}
	return true
}
