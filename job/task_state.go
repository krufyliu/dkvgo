package job

// RunState describe a job state
type TaskState struct {
	Status               string
	FrameAt              int
	PrepareImagesTime    float32
	ComputeFlowTime      float32
	ComputeNovelViewTime float32
	TotalTime            float32
}
