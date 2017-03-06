package protocol

import (
	"github.com/krufyliu/dkvgo/job"
)

type HeartBeatBag struct {
	Todo string
	Echo string
	Task *job.Task
	Report *job.TaskState
}