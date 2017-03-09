package routers

import (
	"github.com/krufyliu/dkvgo/dkvgo-admin/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/", &controllers.MainController{})
}
