package worker

import (
	"flag"
	"testing"
)

func Test_GoodOptions(t *testing.T) {
	args := []string{"-s", "127.0.0.1:9876", "-r", "5", "-w", "60"}
	opt := NewOptions()
	err := opt.fromCmdArgs(args)
	workerFlagSet.PrintDefaults()
	if err != nil {
		t.Error(err.Error())
	}
}

func Test_EmptyAddr(t *testing.T) {
	// protect method `fromCmdArgs` called more than once panic
	workerFlagSet = flag.NewFlagSet("worker", flag.ContinueOnError)
	args := []string{}
	opt := NewOptions()
	err := opt.fromCmdArgs(args)
	if err != errAddrEmpty {
		t.Error("It should return empty addr error")
	}
}

func Test_BadAddr1(t *testing.T) {
	// protect method `fromCmdArgs` called more than once panic
	workerFlagSet = flag.NewFlagSet("worker", flag.ContinueOnError)
	args := []string{"-s", "127.0.0.1:"}
	opt := NewOptions()
	err := opt.fromCmdArgs(args)
	if err != errBadAddr {
		t.Error("It should return bad addr error")
	}
}

func Test_BadAddr2(t *testing.T) {
	// protect method `fromCmdArgs` called more than once panic
	workerFlagSet = flag.NewFlagSet("worker", flag.ContinueOnError)
	args := []string{"-s", "127.0.0.1:abcd"}
	opt := NewOptions()
	err := opt.fromCmdArgs(args)
	if err != errBadAddr {
		t.Error("It should return bad addr error")
	}
}
