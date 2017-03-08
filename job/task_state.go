package job

import (
	"fmt"
)

// TaskState describe a job state
type TaskState struct {
	Status               string
	FrameAt              int
	PrepareImagesTime    float32
	ComputeFlowTime      float32
	ComputeNovelViewTime float32
	TotalTime            float32
}

func (ts TaskState) String() string {
	return fmt.Sprintf("%s-%d(PrepareImages:%f, ComputeFlow:%f, ComputeNovelView:%f, Total:%f)",
		ts.Status, ts.FrameAt, ts.PrepareImagesTime, ts.ComputeFlowTime,
		ts.ComputeNovelViewTime, ts.TotalTime)
}

func (ts TaskState) ShortString() string {
	return fmt.Sprintf("%s-%d", ts.Status, ts.FrameAt)
}
