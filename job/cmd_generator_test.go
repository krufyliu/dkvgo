package job

import (
	"testing"
)

func Test_CmdG(t *testing.T) {
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
	var seg = Task{
		Job: &job,
		Options: &TaskOptions{
			StartFrame: 1200,
			EndFrame:   1299,
		},
	}
	cmdG := NewCmdGeneratorFromTaskSegment(&seg, 0, "/usr/bin", "/etc")
	t.Log(cmdG.GetCmd())
}
