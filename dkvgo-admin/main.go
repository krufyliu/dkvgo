package main

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/krufyliu/dkvgo/dkvgo-admin/models"
	_ "github.com/krufyliu/dkvgo/dkvgo-admin/routers"
)

func init() {
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", "root:mysql1234@/dkvgo?charset=utf8")
	orm.RegisterModel(new(models.User), new(models.Job))
	orm.Debug = true
}

func main() {
	orm.RunCommand()
	beego.Run()
}
