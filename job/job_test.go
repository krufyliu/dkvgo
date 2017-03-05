package job

import (
	"encoding/json"
	"testing"
)

func Test_TaskToJson(t *testing.T) {
	var job = Job{
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
	}
	var out []byte
	var err error
	out, err = json.Marshal(job)
	if err == nil {
		t.Log(string(out))
	}
	job.Map()
	out, err = json.Marshal(job.TaskOpts)
	if err == nil {
		t.Log(string(out))
	}
	var seg = Task{
		Job:    &job,
		Options: job.TaskOpts[0],
	}
	out, err = json.Marshal(seg)
	if err == nil {
		t.Log(string(out))
	}
}
