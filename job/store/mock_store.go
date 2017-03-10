package store

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

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
			Status:            0,
			StartFrame:        1200,
			EndFrame:          1208,
			CameraType:        "AURA",
			Algorithm:         "3D_AURA",
			VideoDir:          "/data/videos/record0001",
			OutputDir:         "/data/output/record0001",
			EnableBottom:      "1",
			EnableTop:         "1",
			Quality:           "8k",
			SaveDebugImg:      "true",
			EanbleColorAdjust: "1",
		},
	}
}

func (store *MockStore) GetJob() *job.Job {
	log.Printf("get task from job store\n")
	var _job = store.job
	store.job = nil
	return _job
}

func (store MockStore) UpdateJob(j *job.Job) bool {
	log.Printf("update job to %d", j.Status)
	return true
}

func (store MockStore) SaveJobState(j *job.Job) bool {
	content, err := json.Marshal(j.TaskOpts)
	if err != nil {
		panic(err)
	}
	log.Printf("save job: %s\n", string(content))
	ioutil.WriteFile(fmt.Sprintf("/tmp/job_%d.state", j.ID), content, 0744)
	return true
}

func (store MockStore) LoadJobState(j *job.Job) bool {
	log.Printf("load job")

	return true
}
