package worker

import (
	"errors"
	"log"
	"net"
	"os/exec"
)

// ErrReachMaxRetryTime means can not connect
const ErrReachMaxRetryTime = errors.New("reach max retry times")

// Worker define worker struct
type Worker struct {
	options    *Options
	connection *net.TCPConn
	cmd        *exec.Cmd
	retry      int
	waitTime   int
}

// NewWorker create a new worker
func NewWorker(opts *Options) *Worker {
	return &Worker{
		options:    opts,
		connection: nil,
		cmd:        nil,
		retry:      0,
		waitTime:   1,
	}
}

func (w *Worker) connect() error {
	conn, err := net.DialTCP("tcp", nil, w.options.schedulerAddr)
	if err != nil {
		return err
	}
	w.retry = 0
	w.waitTime = 1
	w.connection = conn
	return nil
}

func (w *worker) reconnect() error {
	for {
		if w.retry > w.options.maxRetryWaitTime {
			log.Fatal(ErrReachMaxRetryTime)
		}

	}
	return nil
}

func (w *Worker) closeConnection() error {
	if w.connection == nil {
		return nil
	}
	return w.connection.Close()
}

func (w *Worker) stopTask() {
	if w.cmd != nil && w.cmd.ProcessState != nil && !w.cmd.ProcessState.Exited() {
		w.cmd.Process.Kill()
	}
}

func (w *Worker) stop() {
	w.closeConnection()
}

func (w *Worker) nextWaitTime() int {
	w.waitTime += 5
	if w.nextWaitTime > w.options.maxRetryWaitTime {
		w.waitTime = 1
	}
	return w.waitTime
}

func (w *Worker) register() error {
	return nil
}

// RunForever start a forerver goroutine
func (w *Worker) RunForever() {
	for {
		if err := w.connect(); err != nil {
			log.Printf("Error: %s\n", err.Error())
			w.reconnect()
		}
	}
}
