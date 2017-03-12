package main

import (
	"github.com/krufyliu/dkvgo/scheduler"
)

func main() {
	var opts = &scheduler.Options{
		TCPAddr:  ":9876",
		HTTPAddr: "127.0.0.1:9999",
	}
	var sched = scheduler.NewDkvScheduler(opts)
	sched.Main()
}
