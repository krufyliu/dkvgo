package worker

import (
	"bufio"
	"errors"
	"log"
	"net"
	"os/exec"

	"os"

	"io"

	"github.com/krufyliu/dkvgo/job"
)

// ErrReachMaxRetryTime means can not connect
var (
	ErrReachMaxRetryTime = errors.New("reach max retry times")
)

const (
	TaskSubmitAccepted = 1000
	TaskStopAccepted   = 1001
)

// Context define the context of DkvWorker
type Context struct {
	task      *job.Task
	cmd          *exec.Cmd
	state        *job.RunState
	sessionState int
	joined       bool
	sessionID    string
}

// DkvWorker define worker struct
type DkvWorker struct {
	options    *Options
	connection net.Conn
	context    *Context
	retry      int
	waitTime   int
	// task    *job.Task
	// cmd        *exec.Cmd
	// state      *job.RunState
	// joined     bool
	// sessionID  string
}

// NewDkvWorker create a new worker
func NewDkvWorker(opts *Options) *DkvWorker {
	return &DkvWorker{
		options:    opts,
		connection: nil,
		context:    new(Context),
		retry:      0,
		waitTime:   1,
	}
}

func (w *DkvWorker) connect() error {
	conn, err := net.Dial("tcp", w.options.schedulerAddr)
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

func (w *DkvWorker) setSessionState(state int) {
	w.context.sessionState = state
}

func (w *DkvWorker) stopTask() {
	if w.context.cmd != nil && w.context.cmd.ProcessState != nil && !w.context.cmd.ProcessState.Exited() {
		w.context.cmd.Process.Kill()
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

func (w *DkvWorker) getSessionID() string {
	return w.context.sessionID
}

func (w *DkvWorker) runTask(t *job.Task) {
	w.context.task = t
	w.context.cmd = job.NewCmdGeneratorFromTaskSegment(t, 8, "/usr/local/visiondk/bin", "/usr/local/visiondk/setting").GetCmd()
	rd, wd, err := os.Pipe()
	defer rd.Close()
	defer wd.Close()
	if err != nil {
		log.Fatalln(err)
	}
	go w.collectTaskStatus(rd)
	w.context.cmd.Stdout = wd
	w.context.cmd.Stderr = os.Stderr
	err = w.context.cmd.Run()
	if err != nil {
		log.Printf("task %d exited unexpected: %s\n", t.Job.ID, err)
	} else {
		log.Printf("task %d is done\n", t.Job.ID)
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
		w.context.state = state
	}
}

// RunForever start a forerver goroutine
func (w *DkvWorker) RunForever() {
	w.connect()
}
