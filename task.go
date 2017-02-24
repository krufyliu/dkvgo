package dkvgo

import (
	"os/exec"
)

// TaskSplitNum define split level
const TaskSplitNum = 5

// Task define video composition
type Task struct {
	ID        int
	Name      string
	Desc      string
	Priority  int
	Progress  float32
	Status    int
	CreatedAt string
	Options   *TaskOptions
	Segments  []*TaskSegment
}

// Map split t to N small piece
func (t *Task) Map() {
	if len(t.Segments) != 0 {
		return
	}
	var startFrame = t.Options.StartFrame
	var endFrame = t.Options.EndFrame
	var totalFrames = end - startFrame + 1
	var avgFrames = totalFrames / TaskSplitNum
	for i := 0; i < TaskSplitNum; i++ {
		var taskPiece = *t.Options
		taskPiece.StartFrame = startFrame + (i * avgFrames)
		if i == TaskSplitNum-1 {
			taskPiece.EndFrame = endFrame
		} else {
			taskPiece.EndFrame = taskPiece.StartFrame + avgFrames - 1
		}
		t.Segments = append(t.Segments, &taskPiece)
	}
}

// TaskOptions describle video composition parameters
type TaskOptions struct {
	StartFrame        int    `json:"start_frame"`
	EndFrame          int    `json:"end_frame"`
	CameraType        string `json:"camera_type"`
	Algorithm         string `json:"algorithm"` // composition algorithm
	VideoDir          string `json:"video_dir"`
	OutputDir         string `json:"output_dir"`
	EnableBottom      string `json:"enable_bottom"`
	EnableTop         string `json:"enable_top"`
	Quality           string `json:"quality"`
	EanbleColorAdjust string `json:"enable_coloradjust"`
	FrameAt           int    `json:"frame_at"`
}

// TaskSegment describle sub task
type TaskSegment struct {
	TaskID  int
	Options *TaskOptions
}

// GetCmd return a shell cmd for the task
func (t *TaskSegment) GetCmd() *exec.Cmd {

}
