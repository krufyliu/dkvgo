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
	EnableTop        string `orm:"size(10);default(1)"`
	EnableBottom     string `orm:"size(10);default(1)"`
	EnableColorAjust string `orm:"size(10);default(1)"`
	SaveDebugImg     string `orm:"size(10);default(false)"`
	Status           int
	Progress         float32
	Creator          *User `orm:"index;rel(fk)"`
	Operator         *User `orm:"rel(fk)"`
	// CreatorId        int `orm:"index"`
	// OperatorId       int
	CreateAt time.Time `orm:"auto_now_add"`
	UpdateAt time.Time `orm:"auto_now"`
}
