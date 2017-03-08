package scheduler

import (
	"github.com/krufyliu/dkvgo/job"
)

type JobTracker struct {
	ctx     *DkvScheduler
	job     *job.Job
	workers []*Worker
}
