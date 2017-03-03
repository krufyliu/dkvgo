package scheduler

import (
	"github.com/krufyliu/dkvgo/task"
)

type TaskTracer struct {
	ctx     *DkvScheduler
	task    *task.Task
	workers []*Worker
}

func (tr *TaskTracer) Start() {
	tr.ctx.Store.LoadTaskState(tr.task)
	tr.task.Map()
	for _, sopts := range tr.task.SegOpts {
		if sopts.FrameAt == sopts.EndFrame {
			continue
		}
		var taskSeg = &task.TaskSegment {
			Task: tr.task
			Options: sopts
		}
	}
}
