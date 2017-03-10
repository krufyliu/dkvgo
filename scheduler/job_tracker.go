package scheduler

import (
	"log"
	"sync"

	"github.com/krufyliu/dkvgo/job"
	"github.com/krufyliu/dkvgo/job/store"
)

var JobTracker jobTracker

type jobTracker struct {
	sync.RWMutex
	store        store.JobStore
	jobMapping   map[int]map[*job.Task]string
	traceChannel chan *taskSnapShot
}

type taskSnapShot struct {
	task    *job.Task
	state   *job.TaskState
	runAddr string
}

func (tr *jobTracker) trace() {
	for {
		tc := <-tr.traceChannel
		task := tc.task
		_job := task.Job
		if tr.jobMapping[_job.ID] == nil {
			tr.jobMapping[_job.ID] = make(map[*job.Task]string)
		}
		tr.jobMapping[_job.ID][task] = tc.runAddr
		if tc.state != nil {
			tr.handleState(task, tc.state)
		}
	}
}

func (tr *jobTracker) handleState(task *job.Task, state *job.TaskState) {
	_job := task.Job
	var oldState = task.GetState()
	if oldState == nil {
		task.UpdateState(state)
	} else if oldState.FrameAt < state.FrameAt {
		_job.IncFinishFrames(state.FrameAt - oldState.FrameAt)
		task.UpdateState(state)
		tr.store.SaveJobState(_job)
		log.Printf("%s progress: %d/%d/%.2f%%\n",
			_job,
			_job.TotalFrames(),
			_job.GetFinishFrames(),
			_job.CalcProgress())
	}

	switch state.Status {
	case "DONE":
		if _job.TaskDone() {
			if _job.CompareStatusAndSwap(0x05, 0x02, 0x01) {
				tr.store.UpdateJob(_job)
			}
			tr.endTraceJob(_job)
		}
		_job.DecRunning()
	case "STOPPED":
		if _job.DecRunning() == 0 {
			if _job.CompareStatusAndSwap(0x04, 0x03) {
				tr.store.UpdateJob(_job)
			}
			tr.endTraceJob(_job)
		}
	case "FAILED":
		if _job.CompareStatusAndSwap(0x06, 0x02, 0x01) {
			tr.store.UpdateJob(_job)
		}
		if _job.DecRunning() == 0 {
			tr.endTraceJob(_job)
		}
	default:
		if _job.CompareStatusAndSwap(0x02, 0x01) {
			tr.store.UpdateJob(_job)
		}
	}
}

func (tr *jobTracker) TraceWorker(w *Worker) {
	if w.relTask != nil {
		tr.traceChannel <- &taskSnapShot{w.relTask, nil, w.RemoteAddr()}
	}
}

func (tr *jobTracker) TraceWorkerWithState(w *Worker, state *job.TaskState) {
	tr.traceChannel <- &taskSnapShot{w.relTask, state, w.RemoteAddr()}
}

func (tr *jobTracker) endTraceJob(_job *job.Job) {
	delete(tr.jobMapping, _job.ID)
}

func InitJobTracker(store store.JobStore) {
	log.Printf("")
	JobTracker = jobTracker{jobMapping: make(map[int]map[*job.Task]string), traceChannel: make(chan *taskSnapShot, 8), store: store}
	go JobTracker.trace()
}
