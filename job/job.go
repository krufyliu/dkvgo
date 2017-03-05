package job

// TaskSplitNum define split level
const TaskSplitNum = 5

// Job define video composition
type Job struct {
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
	TaskOpts           []*TaskOptions `json:"-"`
}

// Map split t to N small piece
func (t *Job) Map() {
	if len(t.TaskOpts) != 0 {
		return
	}
	var startFrame = t.StartFrame
	var endFrame = t.EndFrame
	var totalFrames = endFrame - startFrame + 1
	var avgFrames = totalFrames / TaskSplitNum
	for i := 0; i < TaskSplitNum; i++ {
		var sOptions = new(TaskOptions)
		sOptions.StartFrame = startFrame + (i * avgFrames)
		sOptions.FrameAt = sOptions.StartFrame
		if i == TaskSplitNum-1 {
			sOptions.EndFrame = endFrame
		} else {
			sOptions.EndFrame = sOptions.StartFrame + avgFrames - 1
		}
		t.TaskOpts = append(t.TaskOpts, sOptions)
	}
}

// TaskOptions describle video composition parameters
type TaskOptions struct {
	StartFrame int `json:"start_frame"`
	EndFrame   int `json:"end_frame"`
	FrameAt    int `json:"frame_at"`
}

// Task describle sub task
type Task struct {
	Job    *Job
	Options *TaskOptions
	Done    bool
}
