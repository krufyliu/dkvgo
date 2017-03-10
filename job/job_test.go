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
		StartFrame:        1000,
		EndFrame:          2500,
		CameraType:        "AURA",
		Algorithm:         "3D_AURA",
		VideoDir:          "/data/videos/record0001",
		OutputDir:         "/data/output/record0001",
		EnableBottom:      "1",
		EnableTop:         "1",
		Quality:           "8k",
		SaveDebugImg:      "true",
		EanbleColorAdjust: "1",
	}
	var out []byte
	var err error
	out, err = json.Marshal(job)
	if err == nil {
		t.Log(string(out))
	}
	job.split()
	out, err = json.Marshal(job.TaskOpts)
	if err == nil {
		t.Log(string(out))
	}
	var seg = Task{
		Job:     &job,
		Options: job.TaskOpts[0],
	}
	out, err = json.Marshal(seg)
	if err == nil {
		t.Log(string(out))
	}
}
