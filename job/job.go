package job

import (
	"fmt"
	"sync"
)

// TaskSplitNum define split level
const TaskSplitNum = 5

// Job define video composition
type Job struct {
	sync.Mutex
	ID                int            `json:"id"`
	Name              string         `json:"name"`
	Priority          int            `json:"priority"`
	Progress          float32        `json:"progress"`
	Status            int            `json:"status"` //0: Queuing, 1: Accepted, 2: Running, 3:stopping, 4:stopped, 5: completed, 6: Failed, 7:Canceled
	StartFrame        int            `json:"start_frame"`
	EndFrame          int            `json:"end_frame"`
	CameraType        string         `json:"camera_type"`
	Algorithm         string         `json:"algorithm"` // composition algorithm
	VideoDir          string         `json:"video_dir"`
	OutputDir         string         `json:"output_dir"`
	EnableBottom      string         `json:"enable_bottom"`
	EnableTop         string         `json:"enable_top"`
	Quality           string         `json:"quality"`
	EanbleColorAdjust string         `json:"enable_coloradjust"`
	TaskOpts          []*TaskOptions `json:"-"`
	numOfCompleteTask int
	finishFrames      int
	numOfTaskRunning  int
}

// map split t to N small piece
func (t *Job) split() {
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

func (t *Job) Init() {
	t.split()
	for _, opt := range t.TaskOpts {
		if opt.FrameAt == opt.EndFrame {
			t.numOfCompleteTask++
			t.finishFrames += opt.EndFrame - opt.FrameAt + 1
		}
	}
}

func (t *Job) GetStatus() int {
	t.Lock()
	defer t.Unlock()
	return t.Status
}

func (t *Job) CompareStatusAndSwap(newStatus int, oldStatus ...int) bool {
	t.Lock()
	defer t.Unlock()
	if len(oldStatus) == 0 {
		t.Status = newStatus
		return true
	}
	for _, value := range oldStatus {
		if t.Status == value {
			t.Status = newStatus
			return true
		}
	}
	return false
}

func (t *Job) TaskDone() bool {
	t.Lock()
	defer t.Unlock()
	t.numOfCompleteTask++
	return t.numOfCompleteTask == len(t.TaskOpts)
}

func (t *Job) IncRunning() int {
	t.Lock()
	defer t.Unlock()
	t.numOfTaskRunning++
	return t.numOfTaskRunning
}

func (t *Job) DecRunning() int {
	t.Lock()
	defer t.Unlock()
	if t.numOfTaskRunning == 0 {
		panic("wrong job running")
	}
	t.numOfTaskRunning--
	return t.numOfTaskRunning
}

func (t *Job) IncFinishFrames(count int) int {
	t.Lock()
	defer t.Unlock()
	t.finishFrames += count
	return t.finishFrames
}

func (t *Job) HasCompleted() bool {
	t.Lock()
	defer t.Unlock()
	return len(t.TaskOpts) != 0 && len(t.TaskOpts) == t.numOfCompleteTask
}

// TaskOptions describle video composition parameters
type TaskOptions struct {
	StartFrame int `json:"start_frame"`
	EndFrame   int `json:"end_frame"`
	FrameAt    int `json:"frame_at"`
}

// Task describle sub task
type Task struct {
	Job     *Job
	Options *TaskOptions
	state   *TaskState
}

func (task Task) String() string {
	return fmt.Sprintf("job[%d](%d-%d:%d)", task.Job.ID, task.Options.StartFrame, task.Options.EndFrame, task.Options.FrameAt)
}

func (task *Task) UpdateState(state *TaskState) {
	task.state = state
	task.Options.FrameAt = state.FrameAt
}

func (task *Task) GetState() *TaskState {
	return task.state
}
