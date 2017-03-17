package models

import (
	"time"
)
// JobState model
type JobState struct {
	Id       int       `orm:"pk;auto"`
	Job *Job `orm:"index;rel(one)"`
	Content  string `orm:"type(text)"`
	CreateAt time.Time `orm:"auto_now_add"`
	UpdateAt time.Time `orm:"auto_now"`
}