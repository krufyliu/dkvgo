package scheduler

import (
	"github.com/krufyliu/dkvgo/job"
)

type JobTracer struct {
	ctx     *DkvScheduler
	job    *job.Job
	workers []*Worker
}

func (tr *JobTracer) Start() {
	tr.ctx.Store.LoadJobState(tr.job)
	tr.job.Map()
	for _, sopts := range tr.job.TaskOpts {
		if sopts.FrameAt == sopts.EndFrame {
			continue
		}
		var task = &job.Task {
			Job: tr.job,
			Options: sopts,
		}
		tr.submit(task)
	}
}

func (tr *JobTracer) submit(task *job.Task) {

}
