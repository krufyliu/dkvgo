package scheduler

import (
	"flag"
	"os"
)

type _options struct {
	TCPAddr  string
	HTTPAddr string
	DBType   string
	DBAddr   string
}

var Options *_options

var tcpAddrHelp = "Listen addr"
var httpAddrHelp = "Http API addr"
var dbTypeHelp = "database type"
var dbAddrHelp = "database dsn"

var _flagSet = flag.NewFlagSet("sched", flag.ContinueOnError)

func _fromCmdArgs(args []string) error {
	if _flagSet.Parsed() {
		panic("This method could only be called once!")
	}
	_flagSet.StringVar(&Options.TCPAddr, "l", ":9876", tcpAddrHelp)
	_flagSet.StringVar(&Options.HTTPAddr, "t", "localhost:9999", httpAddrHelp)
	_flagSet.StringVar(&Options.DBType, "db-type", "mysql", dbTypeHelp)
	_flagSet.StringVar(&Options.DBAddr, "dsn", "", dbAddrHelp)
	return _flagSet.Parse(args)
}

func ParseCmdArgs() {
	Options = &_options{
		TCPAddr:  ":9876",
		HTTPAddr: "localhost:9999",
		DBType:   "mysql",
	}
	if err := _fromCmdArgs(os.Args[1:]); err != nil {
		os.Exit(2)
	}
	if Options.DBAddr == "" {
		os.Stderr.WriteString("dsn must set\n")
		_flagSet.PrintDefaults()
		os.Exit(1)
	}
}
