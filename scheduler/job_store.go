package scheduler

import (
	"github.com/krufyliu/dkvgo/job"
)

// TaskStore define task store and get interface
type TaskStore interface {
	GetJob() *job.Job

	UpdateJob(*job.Job) bool

	SaveJobState(*job.Job) bool

	LoadJobState(*job.Job) bool
}
