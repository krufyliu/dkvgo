package tracker

import (
	"log"
	"sync"

	"time"

	"github.com/krufyliu/dkvgo/job"
	"github.com/krufyliu/dkvgo/job/store"
)

// JobTracker trace job's lifecycle
var JobTracker jobTracker

type jobTracker struct {
	sync.Mutex
	store        store.JobStore
	jobMapping   map[int]map[*job.Task]string
	id2Job       map[int]*job.Job
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
		tr.Lock()
		if tr.jobMapping[_job.ID] == nil {
			log.Printf("begin tracking %s", _job)
			tr.jobMapping[_job.ID] = make(map[*job.Task]string)
			tr.id2Job[_job.ID] = _job
		}
		tr.Unlock()
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
		if time.Since(_job.LastRecord) >= 30*time.Second {
			log.Printf("recording %s, current progress:%.2f", _job, _job.CalcProgress())
			tr.store.UpdateJob(_job)
			tr.store.SaveJobState(_job)
			_job.LastRecord = time.Now()
		}
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
	if state.Status != "RUNNING" {
		tr.store.SaveJobState(_job)
	}
}

func (tr *jobTracker) TraceTask(t *job.Task, addr string) {
	if t != nil {
		tr.traceChannel <- &taskSnapShot{t, nil, addr}
	}
}

func (tr *jobTracker) TraceTaskWithState(t *job.Task, addr string, state *job.TaskState) {
	tr.traceChannel <- &taskSnapShot{t, state, addr}
}

func (tr *jobTracker) endTraceJob(_job *job.Job) {
	tr.Lock()
	delete(tr.jobMapping, _job.ID)
	delete(tr.id2Job, _job.ID)
	tr.Unlock()
	log.Printf("end tracking %s", _job)
}

func (tr *jobTracker) StopJobById(jobID int) bool {
	tr.Lock()
	defer tr.Unlock()
	if _job, ok := tr.id2Job[jobID]; ok {
		log.Printf("try to stop %s", _job)
		if _job.CompareStatusAndSwap(0x3, 0x1, 0x2) {
			return true
		}
	}
	return false
}

func TraceTask(t *job.Task, addr string) {
	JobTracker.TraceTask(t, addr)
}

func TraceTaskWithState(t *job.Task, addr string, state *job.TaskState) {
	JobTracker.TraceTaskWithState(t, addr, state)
}

func StopJobById(jobId int) bool {
	return JobTracker.StopJobById(jobId)
}

// InitWithStore must be called before it is used
func InitWithStore(store store.JobStore) {
	log.Printf("init job tracker")
	JobTracker = jobTracker{
		jobMapping:   make(map[int]map[*job.Task]string),
		id2Job:       make(map[int]*job.Job),
		traceChannel: make(chan *taskSnapShot, 8),
		store:        store,
	}
	go JobTracker.trace()
}
