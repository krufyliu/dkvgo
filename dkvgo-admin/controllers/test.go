package controllers

import (
	"github.com/krufyliu/dkvgo/dkvgo-admin/utils"
)

type TestController struct {
	BaseController
}

func (this *TestController) Md5() {
	this.DataJsonResponse(utils.Md5("visiondk"))	
}