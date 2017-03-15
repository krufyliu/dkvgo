package main

import (
	"github.com/krufyliu/dkvgo/scheduler"
)

func main() {
	var opts = &scheduler.Options{
		TCPAddr:  ":9876",
		HTTPAddr: "127.0.0.1:9999",
		DBType:   "mysql",
		DBAddr:   "root:mysql1234@/dkvgo?charset=utf8",
	}
	var sched = scheduler.NewDkvScheduler(opts)
	sched.Main()
}
