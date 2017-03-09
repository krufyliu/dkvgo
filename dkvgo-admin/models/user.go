package models

import "time"
import "github.com/krufyliu/dkvgo/dkvgo-admin/utils"

// User model
type User struct {
	Id       int       `orm:"pk;auto"`
	Username string    `orm:"unique;size(20)"`
	Email    string    `orm:"unique;size(30)"`
	Password string    `orm:"size(50)" json:"-"`
	CreateAt time.Time `orm:"auto_now_add"`
	UpdateAt time.Time `orm:"auto_now"`
	Jobs     []*Job    `orm:"reverse(many)"`
}

func UserLogin(email, password string) *User, error {
	var md5Password = utils.Md5(password)
}