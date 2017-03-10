package controllers

import (
	"github.com/astaxie/beego"
	"github.com/krufyliu/dkvgo/dkvgo-admin/utils"
)

type LoginController struct {
	beego.Controller
}

func (this LoginController) Post() {
	email := this.GetString("email")
	password := this.GetString("password")
	if email == "" {
		this.Data["json"] = utils.ErrorMap("邮箱不能为空")
		this.ServeJSON()
		return
	}
	if password == "" {
		this.Data["json"] = utils.ErrorMap("密码不能为空")
		this.ServeJSON()
		return
	}
}
