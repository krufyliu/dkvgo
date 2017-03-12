package controllers

import (
	"github.com/krufyliu/dkvgo/dkvgo-admin/models"
	"github.com/krufyliu/dkvgo/dkvgo-admin/services"
)

type UsersController struct {
	BaseController
}

func (this *UsersController) Get() {
	var users []*models.User
	_, err := services.UserService.GetUserList(1, 10).All(&users)
	this.CheckError(err)
	this.DataJsonResponse(users)
}

