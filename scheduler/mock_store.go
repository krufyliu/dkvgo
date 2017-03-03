package scheduler

import (
	"github.com/krufyliu/dkvgo/task"
)

type MockStore struct {
	tk *task.Task
}

func NewMockStore() *MockStore {
	return &MockStore{
		tk: & task.Task{
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
	}
}

func (store MockStore) GetTask() *task.Task {
	return store.tk
}

func (store MockStore) UpdateTask(t *task.Task) bool {
	return true
}

func (store MockStore) SaveTaskState(t *task.Task) bool {
	return true
}

func (store MockStore) SaveTaskState(t *task.Task) bool {
	return true
}