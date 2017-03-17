package controllers

import (
	"time"

	"github.com/krufyliu/dkvgo/dkvgo-admin/services"
)

type AuthController struct {
	BaseController
}

func (this *AuthController) Get() {
	this.DataJsonResponse(this.LoginUser(), "user")
}

func (this *AuthController) Post() {
	email := this.GetString("Email")
	password := this.GetString("Password")
	if email == "" {
		this.ShowErrorMsg("邮箱不能为空")
	}
	if password == "" {
		this.ShowErrorMsg("密码不能为空")
	}
	user, err := services.UserService.GetUserByEmail(email)
	this.CheckError(err)
	if !user.ValidataPassword(password) {
		this.ShowErrorMsg("邮箱不存在或者密码不正确")
	}
	this.SetSession("userId", user.Id)
	user.LastLoginIp = this.GetClientIP()
	user.LastLoginTime = time.Now()
	err = services.UserService.UpdateUser(user, "LastLoginIp", "LastLoginTime")
	this.CheckError(err)
	this.DataJsonResponse(user, "user")
}

func (this *AuthController) Delete() {
	this.DelSession("userId")
	this.ShowSuccessMsg("注销成功")
}
