package worker

import (
	"bufio"
	"errors"
	"log"
	"net"
	"os/exec"

	"encoding/json"

	"os"

	"github.com/krufyliu/dkvgo"
	"github.com/krufyliu/dkvgo/protocol"
)

// ErrReachMaxRetryTime means can not connect
var (
	ErrReachMaxRetryTime = errors.New("reach max retry times")
	ErrVersionNoMatch    = errors.New("version does not match")
	ErrResponseStatus    = errors.New("bad reponse status")
)

// Worker define worker struct
type Worker struct {
	options    *Options
	connection *net.TCPConn
	reader     *bufio.Reader
	taskSeg    *dkvgo.TaskSegment
	cmd        *exec.Cmd
	sessionID  string
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

func (w *Worker) reconnect() error {
	for {
		if w.retry > w.options.maxRetryWaitTime {
			log.Fatal(ErrReachMaxRetryTime)
			break
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
	if w.waitTime > w.options.maxRetryWaitTime {
		w.waitTime = 1
	}
	return w.waitTime
}

func (w *Worker) bufferReader() *bufio.Reader {
	if w.reader != nil {
		return w.reader
	}
	return bufio.NewReader(w.connection)
}

func (w *Worker) register() error {
	var req *protocol.DkvRequst
	if w.sessionID != "" {
		req = protocol.NewDkvRequest("REGISTER", []string{w.sessionID})
	} else {
		req = protocol.NewDkvRequest("REGISTER", []string{})
	}
	message, err := req.Dumps()
	if err != nil {
		return nil
	}
	_, err = w.connection.Write(message)
	if err != nil {
		return nil
	}
	res, err := w.expectResponse()
	if err != nil {
		return err
	}
	if res.Version != "V1.0" {
		return ErrVersionNoMatch
	}
	if res.Status != "OK" {
		return ErrResponseStatus
	}
	if err := res.PullContent(w.bufferReader()); err != nil {
		return err
	}
	return nil
}

func (w *Worker) expectResponse() (*protocol.DkvResponse, error) {
	var res = &protocol.DkvResponse{}
	if err := res.PullHeader(w.bufferReader()); err != nil {
		return nil, err
	}
	return res, nil
}

func (w *Worker) acceptRequest() (*protocol.DkvRequst, error) {
	var req = &protocol.DkvRequst
	if err := req.PullHeader(w.bufferReader()); err != nil {
		return nil, error
	}
	return req, nil
}

func (w *Worker) seedResponse(res *protocol.DkvResponse) error {
	message, err := res.Dumps()
	if err != nil {
		return err
	}
	_, err := w.connection.Write(message)
	return err
}

func (w *Worker) handleStartTask(req *protocol.DkvRequst) {
	var taskSegment dkvgo.TaskSegment
	if err := json.Unmarshal(req.Payload, &taskSegment); err != nil {
		log.Fatalln(err)
	}
	w.seedResponse(protocol.NewDkvResponse("OK"))
	cmdgen := dkvgo.NewCmdGeneratorFromTaskSegment(taskSegment, 8, "/usr/local/visiondk/bin", "/usr/local/visiondk/setting")
	w.taskSeg = &taskSegment
	w.cmd = cmdgen.GetCmd()
}

func (w *Worker) runTask() {
	r, w, err := os.Pipe()
	if err != nil {
		log.Fatalln(err)
	}
	w.cmd.Stdout = w
	w.cmd.Stderr = nil
	err = w.cmd.Run()
	if err != nil {
		if value, ok := err.(*ExitError); ok {
			log.Println(value)
		}
		log.Fatalln(err)
	}
}

// RunForever start a forerver goroutine
func (w *Worker) RunForever() {
	w.connect()
	w.register()
	for {

	}
}
