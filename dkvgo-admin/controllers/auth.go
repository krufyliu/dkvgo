package controllers

import (
	"time"

	"github.com/krufyliu/dkvgo/dkvgo-admin/services"
)

type AuthController struct {
	BaseController
}

func (this *AuthController) Get() {
	if !this.IsLogin() {
		this.Abort(401)
	}
	this.DataJsonResponse(this.LoginUser(), "user")
}

func (this *AuthController) Post() {
	if this.IsLogin() {
		this.ShowErrorMsg("当前已处于登录状态")
	}
	email := this.GetString("email")
	password := this.GetString("password")
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
	this.SetSession("user_id", user.Id)
	user.LastLoginIp = this.GetClientIP()
	user.LastLoginTime = time.Now()
	err = services.UserService.UpdateUser(user, "LastLoginIp", "LastLoginTime")
	this.CheckError(err)
	this.DataJsonResponse(user)
}

func (this *AuthController) Delete() {
	if !this.IsLogin() {
		this.ShowErrorMsg("当前处于登录状态")
	}
	this.DelSession("user_id")
	this.ShowSuccessMsg("注销成功")
}
