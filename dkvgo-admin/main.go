package main

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/krufyliu/dkvgo/dkvgo-admin/routers"
	"github.com/krufyliu/dkvgo/dkvgo-admin/services"
)

func main() {
	services.Init()
	orm.RunCommand()
	beego.SetStaticPath("/assets", "public/assets")
	beego.Run()
}
