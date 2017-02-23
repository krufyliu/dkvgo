package dkvgo

import (
	"os/exec"
)

// Task define video composition
type Task struct {
	ID           int
	Name         string
	Desc         string
	Priority     int
	Progress     float32
	Status       int
	CreatedAt    string
	Options      *TaskOptions
	TaskSegments []*T
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
}

// TaskSegment describle sub task
type TaskSegment struct {
	TaskID  int
	Options *TaskOptions
}

// GetCmd return a shell cmd for the task
func (t *TaskSegment) GetCmd() *exec.Cmd {

}
