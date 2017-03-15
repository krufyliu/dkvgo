package routers

import (
	"github.com/astaxie/beego"
	"github.com/krufyliu/dkvgo/dkvgo-admin/controllers"
)

func init() {
	apiNs := beego.NewNamespace("/api",
		beego.NSRouter("/auth", &controllers.AuthController{}),
		beego.NSRouter("/users", &controllers.UsersController{}),
		beego.NSRouter("/jobs", &controllers.JobsController{}),
		beego.NSRouter("/jobs/:id:int/action/stop", &controllers.JobsController{}, "post:Stop"),
		beego.NSRouter("/jobs/:id:int/action/resume", &controllers.JobsController{}, "post:Resume"),
	)
	beego.AddNamespace(apiNs)
	beego.AutoRouter(&controllers.TestController{})
	beego.Router("*", &controllers.MainController{})
}
