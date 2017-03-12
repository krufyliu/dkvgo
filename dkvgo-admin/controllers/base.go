package controllers

import (
	"github.com/astaxie/beego"
)

const (
	MSG_OK       = true  // ajax输出错误码，成功
	MSG_ERR      = false // 错误
)

type BaseController struct {
	beego.Controller
}

func (this *BaseController) ShouldReturnJson() bool {
	return this.IsAjax() || this.Ctx.Input.AcceptsJSON()
}

func (this *BaseController) ShowMsg(msg string, success bool) {
	out := make(map[string]interface{})
	out["success"] = success
	out["message"] = msg
	this.JsonResponse(out)
}

func (this *BaseController) GetClientIP() string {
	if p := this.Ctx.Input.Proxy(); len (p) > 0 {
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
	out["status"] = MSG_OK
	out["message"] = "success"
	out[fieldname] = data
	this.JsonResponse(out)
}

func (this *BaseController) ErrorJsonResponse(msg string, detail interface{}) {
	out := make(map[string]interface{})
	out["status"] = MSG_ERR
	out["message"] = msg
	out["errors"] = detail
	this.JsonResponse(out)
}

func (this *BaseController) CheckError(err error) {
	if err == nil {
		return
	}
	this.ShowMsg(err.Error(), MSG_ERR)
}