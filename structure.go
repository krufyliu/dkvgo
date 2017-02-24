package dkvgo

// FrameState describe image handle state
type FrameState struct {
	FrameAt              int
	PrepareImagesTime    float32
	ComputeFlowTime      float32
	ComputeNovelViewTime float32
	TotalTime            float32
}
