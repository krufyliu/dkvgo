package dkvgo

// TaskState describe a job state
type TaskState struct {
	Running              bool
	FrameAt              int
	PrepareImagesTime    float32
	ComputeFlowTime      float32
	ComputeNovelViewTime float32
	TotalTime            float32
}

type PeerInfo struct {
	id         string
	lastUpdate int64
	CmdAddr    string
	Version    string
}
