package main

import (
	"github.com/krufyliu/dkvgo/scheduler"
)

func main() {
	scheduler.ParseCmdArgs()
	scheduler.Run()
}
