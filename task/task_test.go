package task

import (
	"encoding/json"
	"testing"
)

func Test_TaskToJson(t *testing.T) {
	var task = Task{
		ID:                1,
		Name:              "test",
		Priority:          128,
		Progress:          88.5,
		Status:            1,
		CreatedAt:         141241212312,
		StartFrame:        1000,
		EndFrame:          2500,
		CameraType:        "GOPRO",
		Algorithm:         "FACEBOOK_3D",
		VideoDir:          "/data/videos/record0001",
		OutputDir:         "/data/output/record0001",
		EnableBottom:      "1",
		EnableTop:         "1",
		Quality:           "8k",
		EanbleColorAdjust: "1",
		FrameAt:           1200,
	}
	var out []byte
	var err error
	out, err = json.Marshal(task)
	if err == nil {
		t.Log(string(out))
	}
	task.Map()
	out, err = json.Marshal(task.SegOpts)
	if err == nil {
		t.Log(string(out))
	}
	var seg = TaskSegment{
		Task:    &task,
		Options: task.SegOpts[0],
	}
	out, err = json.Marshal(seg)
	if err == nil {
		t.Log(string(out))
	}
}
