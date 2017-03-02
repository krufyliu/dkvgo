package scheduler

import (
	"github.com/krufyliu/dkvgo/task"
)

// TaskStore define task store and get interface
type TaskStore interface {
	GetTask() *task.Task

	UpdateTask(t *task.Task) bool

	SaveTaskState(t *task.Task) bool
}
