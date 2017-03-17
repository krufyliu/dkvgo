package services

import (
	"errors"
	"github.com/astaxie/beego/orm"
	"github.com/krufyliu/dkvgo/dkvgo-admin/models"
)

type userService struct {} 

func (this *userService) GetTotal() (int64, error) {
	return o.QueryTable(&models.User{}).Count()
}

func (this *userService) GetPage(current, pageSize int) (*Page, error) {
	total, err := this.GetTotal()
	if err != nil {
		return nil, err
	}
	return &Page {Total: total, Current: current, PageSize: pageSize}, nil
}

func (this *userService) GetUserById(userId int) (*models.User, error) {
	user := &models.User{Id: userId}
	err := o.Read(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (this *userService) GetUserByEmail(email string) (*models.User, error) {
	user := &models.User{}
	user.Email = email
	err := o.Read(user, "Email")
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (this *userService) GetUserList(page, pageSize int) orm.QuerySeter {
	offset := (page-1) *pageSize
	if offset < 0 {
		offset = 0
	}
	qs := o.QueryTable(&models.User{}).Limit(pageSize, offset)
	return qs
}

func (this *userService) UpdateUser(user *models.User, fields ...string) error{
	if len(fields) < 1 {
		return errors.New("更新字段不能为空")
	}
	_, err := o.Update(user, fields...)
	return err
}
