package models

import (
	"time"

	"github.com/krufyliu/dkvgo/dkvgo-admin/utils"
)

// User model
type User struct {
	Id            int       `orm:"pk;auto"`
	Username      string    `orm:"unique;size(20)"`
	Email         string    `orm:"unique;size(30)"`
	Password      string    `orm:"size(50)" json:"-"`
	LastLoginIp   string    `orm:"size(30)"`
	LastLoginTime time.Time `orm:"null"`
	CreateAt      time.Time `orm:"auto_now_add"`
	UpdateAt      time.Time `orm:"auto_now"`
	Jobs          []*Job    `orm:"reverse(many)"`
}

func (this *User) ValidataPassword(password string) bool {
	hashPassword := utils.Md5(password)
	return this.Password == hashPassword
}

func (this *User) IsAdmin() bool {
	return this.Id == 1
}
