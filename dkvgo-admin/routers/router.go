package routers

import (
	"github.com/krufyliu/dkvgo/dkvgo-admin/controllers"
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/auth", &controllers.AuthController{})
	beego.Router("/users", &controllers.UsersController{})
	beego.Router("/jobs", &controllers.JobsController{})
	beego.AutoRouter(&controllers.TestController{})
    beego.Router("/", &controllers.MainController{})
}
