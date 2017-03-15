package worker

import (
	"bufio"
	"log"
	"net"
	"os/exec"
	"time"

	"os"

	"io"

	"github.com/krufyliu/dkvgo/job"
)

// Context define the ctx of DkvWorker
type Context struct {
	task      *job.Task
	cmd       *exec.Cmd
	state     *job.TaskState
	forceStop bool
}

// DkvWorker define worker struct
type DkvWorker struct {
	options    *Options
	connection net.Conn
	ctx        *Context
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
		ctx:        nil,
		retry:      0,
		waitTime:   5,
	}
}

func (w *DkvWorker) connect() error {
	conn, err := net.Dial("tcp", w.options.schedulerAddr)
	if err != nil {
		log.Printf("Error: %s\n", err)
		return err
	}
	w.retry = 0
	w.waitTime = 5
	w.connection = conn
	return nil
}

func (w *DkvWorker) tryToConnect() {
	for {
		if err := w.connect(); err != nil {
			w.retry++
		} else {
			break
		}
		if w.options.maxRetry != 0 && w.retry > w.options.maxRetry {
			log.Fatalf("Error: reach max retry connect times: %d\n", w.options.maxRetry)
		}
		var sleepTime = w.nextWaitTime()
		log.Printf("after %ds to reconnect, reconnect times: %d\n", sleepTime, w.retry-1)
		time.Sleep(time.Duration(sleepTime) * time.Second)
	}
}

func (w *DkvWorker) nextWaitTime() int {
	var waitTime = w.waitTime
	w.waitTime += 5
	if w.waitTime > w.options.maxRetryWaitTime {
		w.waitTime = 5
	}
	return waitTime
}

func (w *DkvWorker) closeConnection() error {
	if w.connection == nil {
		return nil
	}
	w.joined = false
	w.sessionID = ""
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
	if w.ctx != nil && w.ctx.cmd != nil && w.ctx.cmd.Process != nil {
		log.Printf("stop task ...\n")
		w.ctx.cmd.Process.Kill()
		log.Printf("task stopped\n")
	}
}

func (w *DkvWorker) stop() {
	w.closeConnection()
	w.stopTask()
	w.clearCtx()
}

func (w *DkvWorker) initCtx(t *job.Task) {
	var cg = job.NewCmdGeneratorFromTaskSegment(t, 8, "/usr/local/visiondk/bin", "/usr/local/visiondk/setting")
	var ctx = &Context{
		state: &job.TaskState{FrameAt: t.Options.FrameAt},
		cmd:   cg.GetCmd(),
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
	w.ctx.cmd.Stderr = wd
	w.ctx.state.Status = "RUNNING"
	log.Printf("run shell task %+v\n", w.ctx.cmd.Args)
	var currentSession = w.sessionID
	err = w.ctx.cmd.Run()
	if currentSession != w.sessionID {
		return
	}
	//try to make full collection
	time.Sleep(1 * time.Second)
	if err != nil {
		if w.ctx.forceStop {
			w.ctx.state.Status = "STOPPED"
		} else {
			w.ctx.state.Status = "FAILED"
		}
		log.Printf("%s exited unexpected: %s\n", t, err)
	} else {
		w.ctx.state.Status = "DONE"
		log.Printf("%s is done\n", t)
	}
}

func (w *DkvWorker) collectTaskStatus(r io.Reader) {
	reader := bufio.NewReader(r)
	for {
		state, err := matchState(reader)
		if err != nil && err != io.EOF {
			log.Fatalf("Error: %s\n", err)
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
	for {
		w.tryToConnect()
		var pl = ProtocolLoop{ctx: w}
		if err := pl.IOLoop(w.connection); err != nil {
			log.Printf("Error: %s------------------\n", err)
			w.stop()
		}
	}
}
