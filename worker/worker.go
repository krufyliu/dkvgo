package worker

import (
	"errors"
	"log"
	"net"
	"os/exec"

	"os"

	"github.com/krufyliu/dkvgo"
)

// ErrReachMaxRetryTime means can not connect
var (
	ErrReachMaxRetryTime = errors.New("reach max retry times")
	ErrVersionNoMatch    = errors.New("version does not match")
	ErrResponseStatus    = errors.New("bad reponse status")
)

// DkvWorker define worker struct
type DkvWorker struct {
	options    *Options
	connection *net.TCPConn
	taskSeg    *dkvgo.TaskSegment
	cmd        *exec.Cmd
	sessionID  string
	retry      int
	waitTime   int
}

// NewDkvWorker create a new worker
func NewDkvWorker(opts *Options) *DkvWorker {
	return &DkvWorker{
		options:    opts,
		connection: nil,
		cmd:        nil,
		retry:      0,
		waitTime:   1,
	}
}

func (w *DkvWorker) connect() error {
	conn, err := net.DialTCP("tcp", nil, w.options.schedulerAddr)
	if err != nil {
		return err
	}
	w.retry = 0
	w.waitTime = 1
	w.connection = conn
	return nil
}

func (w *DkvWorker) reconnect() error {
	for {
		if w.retry > w.options.maxRetryWaitTime {
			log.Fatal(ErrReachMaxRetryTime)
			break
		}

	}
	return nil
}

func (w *DkvWorker) closeConnection() error {
	if w.connection == nil {
		return nil
	}
	return w.connection.Close()
}

func (w *DkvWorker) stopTask() {
	if w.cmd != nil && w.cmd.ProcessState != nil && !w.cmd.ProcessState.Exited() {
		w.cmd.Process.Kill()
	}
}

func (w *DkvWorker) stop() {
	w.stopTask()
	w.closeConnection()
}

func (w *DkvWorker) nextWaitTime() int {
	w.waitTime += 5
	if w.waitTime > w.options.maxRetryWaitTime {
		w.waitTime = 1
	}
	return w.waitTime
}

func (w *DkvWorker) runTask() {
	_, wd, err := os.Pipe()
	if err != nil {
		log.Fatalln(err)
	}

	w.cmd.Stdout = wd
	w.cmd.Stderr = nil
	err = w.cmd.Run()
	if err != nil {
		if value, ok := err.(*exec.ExitError); ok {
			log.Println(value)
		}
		log.Fatalln(err)
	}
}

// RunForever start a forerver goroutine
func (w *DkvWorker) RunForever() {
}
