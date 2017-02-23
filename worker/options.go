package worker

import (
	"errors"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
)

// DefaultSchedulerPort present default scheduler port
const DefaultSchedulerPort = 7654

// Options define worker runtime option
type Options struct {
	schedulerAddr    *net.TCPAddr
	maxRetry         int
	maxRetryWaitTime int
}

// NewOptions create default worker options
func NewOptions() *Options {
	return &Options{
		maxRetry:         0,
		maxRetryWaitTime: 90,
	}
}

var schedulerAddHelp = "Schduler tcp addr, non empty, eg 127.0.0.1:8765. Defualt port 7654."
var maxRetryHelp = "Max retry times when lost connection."
var maxRetryWaitTimeHelp = "Max wait time before next connection."

var errAddrEmpty = errors.New("scheduler tcp addr")
var errBadAddr = errors.New("bad scheduler tcp addr")

var workerFlagSet = flag.NewFlagSet("worker", flag.ContinueOnError)

func (opt *Options) fromCmdArgs(args []string) error {
	if workerFlagSet.Parsed() {
		panic("This method could only be called once!")
	}
	var schedulerAddr string
	var maxRetry int
	var maxRetryWaitTime int

	workerFlagSet.StringVar(&schedulerAddr, "scheduler-addr", "", schedulerAddHelp)
	workerFlagSet.StringVar(&schedulerAddr, "s", "", schedulerAddHelp+" short for scheduler-addr")
	workerFlagSet.IntVar(&maxRetry, "max-retry", 0, maxRetryHelp)
	workerFlagSet.IntVar(&maxRetry, "r", 0, maxRetryHelp+" short for max-retry")
	workerFlagSet.IntVar(&maxRetryWaitTime, "max-retry-wait-time", 90, maxRetryWaitTimeHelp)
	workerFlagSet.IntVar(&maxRetryWaitTime, "w", 90, maxRetryWaitTimeHelp+" short for max-retry-wait-time")

	err := workerFlagSet.Parse(args)
	if err != nil {
		return err
	}
	if schedulerAddr == "" {
		return errAddrEmpty
	}
	if !strings.Contains(schedulerAddr, ":") {
		schedulerAddr = fmt.Sprintf("%s:%d", schedulerAddr, DefaultSchedulerPort)
	}
	tcpAddr, err := net.ResolveTCPAddr("tcp", schedulerAddr)
	if err != nil || tcpAddr.Port == 0 {
		return errBadAddr
	}

	opt.schedulerAddr = tcpAddr
	opt.maxRetry = maxRetry
	opt.maxRetryWaitTime = maxRetryWaitTime
	return nil
}

func (opt *Options) tryFromCmdArgs(args []string) {
	if err := opt.fromCmdArgs(args); err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		workerFlagSet.PrintDefaults()
		os.Exit(2)
	}
}
