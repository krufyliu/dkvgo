package task

// TaskSplitNum define split level
const TaskSplitNum = 5

// Task define video composition
type Task struct {
	ID                int               `json:"id"`
	Name              string            `json:"name"`
	Priority          int               `json:"priority"`
	Progress          float32           `json:"progress"`
	Status            int               `json:"status"`
	CreatedAt         int               `json:"created_at"`
	StartFrame        int               `json:"start_frame"`
	EndFrame          int               `json:"end_frame"`
	CameraType        string            `json:"camera_type"`
	Algorithm         string            `json:"algorithm"` // composition algorithm
	VideoDir          string            `json:"video_dir"`
	OutputDir         string            `json:"output_dir"`
	EnableBottom      string            `json:"enable_bottom"`
	EnableTop         string            `json:"enable_top"`
	Quality           string            `json:"quality"`
	EanbleColorAdjust string            `json:"enable_coloradjust"`
	SegOpts           []*SegmentOptions `json:"-"`
}

// Map split t to N small piece
func (t *Task) Map() {
	if len(t.SegOpts) != 0 {
		return
	}
	var startFrame = t.StartFrame
	var endFrame = t.EndFrame
	var totalFrames = endFrame - startFrame + 1
	var avgFrames = totalFrames / TaskSplitNum
	for i := 0; i < TaskSplitNum; i++ {
		var sOptions = new(SegmentOptions)
		sOptions.StartFrame = startFrame + (i * avgFrames)
		sOptions.FrameAt = sOptions.StartFrame
		if i == TaskSplitNum-1 {
			sOptions.EndFrame = endFrame
		} else {
			sOptions.EndFrame = sOptions.StartFrame + avgFrames - 1
		}
		t.SegOpts = append(t.SegOpts, sOptions)
	}
}

// SegmentOptions describle video composition parameters
type SegmentOptions struct {
	StartFrame int `json:"start_frame"`
	EndFrame   int `json:"end_frame"`
	FrameAt    int `json:"frame_at"`
}

// TaskSegment describle sub task
type TaskSegment struct {
	Task    *Task
	Options *SegmentOptions
}
