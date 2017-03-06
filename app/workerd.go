package main

import (
	"os"

	"github.com/krufyliu/dkvgo/worker"
)

func main() {
	var opts = worker.NewOptions()
	opts.TryFromCmdArgs(os.Args)
	var workerd = worker.NewDkvWorker(opts)
	workerd.Main()
}
