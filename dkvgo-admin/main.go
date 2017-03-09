package main

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/krufyliu/dkvgo/dkvgo-admin/routers"
)

func main() {
	orm.RunCommand()
	beego.Run()
}
