package models

import "time"

// Job model
type Job struct {
	Id               int    `orm:"pk;auto`
	Name             string `orm:"size(50)"`
	VideoDir         string `orm:"size(512)"`
	OutputDir        string `orm:"size(512)"`
	StartFrame       int
	EndFrame         int
	Algorithm        string `orm:"size(20)"`
	Priority         int
	CameraType       string `orm:"size(10)"`
	Quality          string `orm:"size(10)"`
	EnableTop        string `orm:"size(1);default(1)"`
	EnableBottom     string `orm:"size(1);default(1)"`
	EnableColorAdjust string `orm:"size(1);default(1)"`
	SaveDebugImg     string `orm:"size(5);default(false)"`
	Status           int
	Progress         float32
	Creator          *User `orm:"index;rel(fk)"`
	Operator         *User `orm:"rel(fk)"`
	State *JobState `orm:"reverse(one)"`
	// CreatorId        int `orm:"index"`
	// OperatorId       int
	CreateAt time.Time `orm:"auto_now_add"`
	UpdateAt time.Time `orm:"auto_now"`
}
