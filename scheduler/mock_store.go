package scheduler

import (
	"github.com/krufyliu/dkvgo/job"
)

type MockStore struct {
	job *job.Job
}

func NewMockStore() *MockStore {
	return &MockStore{
		job: &job.Job{
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
		},
	}
}

func (store MockStore) GetJob() *job.Job {
	return store.job
}

func (store MockStore) UpdateJob(j *job.Job) bool {
	return true
}

func (store MockStore) SaveJobState(j *job.Job) bool {
	return true
}

func (store MockStore) LoadJobState(j *job.Job) bool {
	return true
}