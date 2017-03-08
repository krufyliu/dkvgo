package store

import (
	"github.com/krufyliu/dkvgo/job"
)

// TaskStore define task store and get interface
type JobStore interface {
	// fetch preparing job
	GetJob() *job.Job

	// update job
	UpdateJob(*job.Job) bool

	// save the state of a job
	SaveJobState(*job.Job) bool

	// get the state of a job
	LoadJobState(*job.Job) bool
}
