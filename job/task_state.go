package job

// RunState describe a job state
type RunState struct {
	Status               int
	FrameAt              int
	PrepareImagesTime    float32
	ComputeFlowTime      float32
	ComputeNovelViewTime float32
	TotalTime            float32
}
