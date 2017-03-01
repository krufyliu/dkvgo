package worker

import (
	"bufio"
	"errors"
	"log"
	"net"
	"os/exec"

	"os"

	"io"

	"github.com/krufyliu/dkvgo/task"
)

// ErrReachMaxRetryTime means can not connect
var (
	ErrReachMaxRetryTime = errors.New("reach max retry times")
)

// Context define the context of DkvWorker
type Context struct {
	taskSeg   *task.TaskSegment
	cmd       *exec.Cmd
	state     *task.RunState
	joined    bool
	sessionID string
}

// DkvWorker define worker struct
type DkvWorker struct {
	options    *Options
	connection *net.TCPConn
	context    *Context
	retry      int
	waitTime   int
	// taskSeg    *task.TaskSegment
	// cmd        *exec.Cmd
	// state      *task.RunState
	// joined     bool
	// sessionID  string
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
	w.closeConnection()
	w.stopTask()
}

func (w *DkvWorker) nextWaitTime() int {
	w.waitTime += 5
	if w.waitTime > w.options.maxRetryWaitTime {
		w.waitTime = 1
	}
	return w.waitTime
}

func (w *DkvWorker) setSessionID(sessionID string) {
	if w.context.sessionID != "" && w.context.sessionID != sessionID {
		w.stopTask()
		w.context.sessionID = sessionID
	}
}

func (w *DkvWorker) runTask(t *task.TaskSegment) {
	w.context.taskSeg = t
	w.context.cmd = task.NewCmdGeneratorFromTaskSegment(t, 8, "/usr/local/visiondk/bin", "/usr/local/visiondk/setting").GetCmd()
	rd, wd, err := os.Pipe()
	defer rd.Close()
	defer wd.Close()
	if err != nil {
		log.Fatalln(err)
	}
	go w.collectTaskStatus(rd)
	w.cmd.Stdout = wd
	w.cmd.Stderr = os.Stderr
	err = w.cmd.Run()
	if err != nil {
		log.Printf("task %d exited unexpected: %s\n", t.Task.ID, value.String())
	} else {
		log.Printf("task %d is done\n", t.Task.ID)
	}
}

func (w *DkvWorker) collectTaskStatus(r io.Reader) {
	reader := bufio.NewReader(r)
	for {
		state, err := matchState(reader)
		if err != nil && err != io.EOF {
			return
		}
		if err == io.EOF {
			break
		}
		w.state = state
	}
}

// RunForever start a forerver goroutine
func (w *DkvWorker) RunForever() {
	w.connect()
}
