package controllers


import (
	"github.com/krufyliu/dkvgo/dkvgo-admin/services"
	"time"
)

type AuthController struct {
	BaseController
}

func (this *AuthController) Post() {
	email := this.GetString("email")
	password := this.GetString("password")
	if email == "" {
		this.ShowMsg("邮箱不能为空", MSG_ERR)
	}
	if password == "" {
		this.ShowMsg("密码不能为空", MSG_ERR)
	}
	user, err := services.UserService.GetUserByEmail(email)
	this.CheckError(err)
	if !user.ValidataPassword(password) {
		this.ShowMsg("邮箱不存在或者密码不正确", MSG_ERR)
	}
	this.SetSession("user_id", user.Id)
	user.LastLoginIp = this.GetClientIP()
	user.LastLoginTime = time.Now()
	err = services.UserService.UpdateUser(user, "LastLoginIp", "LastLoginTime")
	this.CheckError(err)
	this.DataJsonResponse(user)
}

func (this *AuthController) Delete() {
	this.DelSession("user_id")
	this.ShowMsg("注销成功", MSG_OK)
}