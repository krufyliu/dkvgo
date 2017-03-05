package worker

import (
	"bufio"
	"errors"
	"log"
	"net"
	"time"
	"os/exec"

	"os"

	"io"

	"github.com/krufyliu/dkvgo/job"
)

// ErrReachMaxRetryTime means can not connect
var (
	ErrReachMaxRetryTime = errors.New("reach max retry times")
)

// Context define the ctx of DkvWorker
type Context struct {
	task      	 *job.Task
	cmd          *exec.Cmd
	state        *job.TaskState
	forceStop    bool
}

// DkvWorker define worker struct
type DkvWorker struct {
	options    *Options
	connection net.Conn
	ctx    	   *Context
	lastUpdate int64
	retry      int
	waitTime   int
	joined     bool
	sessionID  string
}

// NewDkvWorker create a new worker
func NewDkvWorker(opts *Options) *DkvWorker {
	return &DkvWorker{
		options:    opts,
		connection: nil,
		ctx:    nil,
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

func (w *DkvWorker) clearCtx() {
	w.ctx = nil
}

func (w *DkvWorker) updatePing() {
	w.lastUpdate = time.Now().Unix()
}


func (w *DkvWorker) forceStopTask() {
	w.ctx.forceStop = true
	w.stopTask()
}

func (w *DkvWorker) stopTask() {
	if w.ctx.cmd != nil && w.ctx.cmd.ProcessState != nil && !w.ctx.cmd.ProcessState.Exited() {
		w.ctx.cmd.Process.Kill()
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

func (w *DkvWorker) initCtx(t *job.Task) {
	var ctx = &Context {
		state: new(job.TaskState),
		cmd: job.NewCmdGeneratorFromTaskSegment(t, 8, "/usr/local/visiondk/bin", "/usr/local/visiondk/setting").GetCmd(),
	}
	w.ctx = ctx
}

func (w *DkvWorker) runTask(t *job.Task) {
	w.initCtx(t)
	rd, wd, err := os.Pipe()
	defer rd.Close()
	defer wd.Close()
	if err != nil {
		log.Fatalln(err)
	}
	go w.collectTaskStatus(rd)
	w.ctx.cmd.Stdout = wd
	w.ctx.cmd.Stderr = os.Stderr
	w.ctx.state.Status = "RUNNING"
	err = w.ctx.cmd.Run()
	w.ctx.state.Status = "Done"
	if err != nil {
		if w.ctx.forceStop {
			w.ctx.state.Status = "STOPPED"
		} else {
			w.ctx.state.Status = "Failed"
		}
		log.Printf("task %d exited unexpected: %s\n", t.Job.ID, err)
	} 
	log.Printf("task %d is done\n", t.Job.ID)
	w.clearCtx()
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
		state.Status = w.ctx.state.Status
		w.ctx.state = state
	}
}

func (w *DkvWorker) Main() {
	w.connect()
	var pl = ProtocolLoop{ctx: w}
	pl.IOLoop(w.connection)
}
