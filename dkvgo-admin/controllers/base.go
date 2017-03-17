package controllers

import (
	"github.com/astaxie/beego"
	"github.com/krufyliu/dkvgo/dkvgo-admin/models"
	"github.com/krufyliu/dkvgo/dkvgo-admin/services"
)

const (
	MSG_OK  = true  // ajax输出错误码，成功
	MSG_ERR = false // 错误
)

type BaseController struct {
	beego.Controller
	loginUser *models.User
}

func (this *BaseController) Prepare() {
	if this.Ctx.Input.Method() == "POST" && this.Ctx.Input.URI() == "/api/auth" {
		if this.IsLogin() {
			this.ShowErrorMsg("已经登录过了")
		}
	} else {
		if !this.IsLogin() {
			this.ShowErrorMsg("未登录")
		}
	}
}

func (this *BaseController) LoginUser() *models.User {
	if this.loginUser != nil {
		return this.loginUser
	} else if this.GetSession("userId") != nil {
		loginUser, err := services.UserService.GetUserById(this.GetSession("userId").(int))
		this.CheckError(err)
		this.loginUser = loginUser
		return loginUser
	}
	return nil
}

func (this *BaseController) ShouldReturnJson() bool {
	return this.IsAjax() || this.Ctx.Input.AcceptsJSON()
}

func (this *BaseController) IsLogin() bool {
	return this.LoginUser() != nil
}

func (this *BaseController) ShowMsg(msg string, success bool) {
	out := make(map[string]interface{})
	out["success"] = success
	out["message"] = msg
	this.JsonResponse(out)
}

func (this *BaseController) ShowErrorMsg(msg string) {
	this.ShowMsg(msg, MSG_ERR)
}

func (this *BaseController) ShowSuccessMsg(msg string) {
	this.ShowMsg(msg, MSG_OK)
}

func (this *BaseController) GetClientIP() string {
	if p := this.Ctx.Input.Proxy(); len(p) > 0 {
		return p[0]
	}
	return this.Ctx.Input.IP()
}

func (this *BaseController) JsonResponse(out interface{}) {
	this.Data["json"] = out
	this.ServeJSON()
	this.StopRun()
}

func (this *BaseController) DataJsonResponse(data interface{}, field ...string) {
	fieldname := "data"
	if len(field) > 0 {
		fieldname = field[0]
	}
	out := make(map[string]interface{})
	out["success"] = MSG_OK
	out["message"] = "success"
	out[fieldname] = data
	this.JsonResponse(out)
}

func (this *BaseController) DataJsonResponseWithPage(data interface{}, page interface{}) {
	out := make(map[string]interface{})
	out["success"] = MSG_OK
	out["message"] = "success"
	out["data"] = data
	out["page"] = page
	this.JsonResponse(out)
}

func (this *BaseController) ErrorJsonResponse(msg string, detail interface{}) {
	out := make(map[string]interface{})
	out["success"] = MSG_ERR
	out["message"] = msg
	out["errors"] = detail
	this.JsonResponse(out)
}

func (this *BaseController) CheckError(err error) {
	if err == nil {
		return
	}
	beego.Error(err)
	this.ShowMsg(err.Error(), MSG_ERR)
}
