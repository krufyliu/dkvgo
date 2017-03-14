package controllers

import (
	"github.com/krufyliu/dkvgo/dkvgo-admin/models"
	"github.com/krufyliu/dkvgo/dkvgo-admin/services"
)

type UsersController struct {
	BaseController
}

func (this *UsersController) Get() {
	page, err := this.GetInt("page", 1)
	this.CheckError(err)
	pageSize, err := this.GetInt("size", 10)
	this.CheckError(err)
	var users []*models.User
	qs := services.UserService.GetUserList(page, pageSize)
	_, err = qs.OrderBy("-UpdateAt").All(&users)
	this.CheckError(err)
	pager, err := services.UserService.GetPage(page, pageSize)
	this.DataJsonResponseWithPage(users, pager)
}
