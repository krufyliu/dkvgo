package scheduler

import (
	"bufio"
	"crypto/md5"
	"encoding/json"
	"errors"
	"log"
	"net"
	"time"

	"io"

	"github.com/krufyliu/dkvgo/job"
	"github.com/krufyliu/dkvgo/protocol"
)

type Worker struct {
	ctx        *DkvScheduler
	conn       net.Conn
	reader     *bufio.Reader
	remoteAddr string
	lastUpdate int64
	relTask    *job.Task
}

func (worker *Worker) Attach(ts *job.Task) {
	worker.relTask = ts
}

func (worker *Worker) Dettach() *job.Task {
	var ts = worker.relTask
	worker.relTask = nil
	return ts
}

func (worker *Worker) IOLoop(conn net.Conn) {
	worker.conn = conn
	worker.reader = bufio.NewReader(conn)
	worker.remoteAddr = conn.RemoteAddr().String()
	if err := worker.identity(); err != nil {
		log.Printf("Error: %s\n", err)
		worker.conn.Close()
		return
	}
	worker.runLoop()
}

func (worker *Worker) identity() error {
	log.Printf("waiting for worker[%s] identity\n", worker.conn.RemoteAddr().String())
	var deadline = time.Now().Add(2 * time.Second)
	worker.conn.SetReadDeadline(deadline)
	var pack = new(protocol.Package)
	if err := pack.Unmarshal(worker.reader); err != nil {
		return err
	}
	if pack.Directive != 0x01 {
		return errors.New("not a register request")
	}
	hash := md5.New()
	hash.Write([]byte(worker.remoteAddr))
	hash.Write([]byte(time.Now().String()))
	pack = protocol.NewPackageWithPayload(0x01, hash.Sum(nil))
	if err := worker.sendPackage(pack); err != nil {
		return err
	}
	return nil
}

func (worker *Worker) runLoop() {
	for {
		err := worker.makeOnePullResponse()
		if err != nil {
			if err == io.EOF {
				worker.conn.Close()
				if worker.relTask != nil {
					worker.ctx.TaskPool.PushFront(worker.relTask)
				}
				break
			}
			log.Printf("Error: %s\n", err)
		}
	}
}

func (worker *Worker) makeOnePullResponse() error {
	pack, err := worker.receivePackage()
	if err != nil {
		return err
	}
	worker.handleHeartBeatRequest(pack)
	return nil
}

func (worker *Worker) handleHeartBeatRequest(pack *protocol.Package) error {
	if pack.Directive != 0x02 && pack.WithPack != 1 {
		return errors.New("bad heartbeat request")
	}
	var bag = new(protocol.HeartBeatBag)
	err := json.Unmarshal(pack.Payload, bag)
	if err != nil {
		return nil
	}
	switch bag.Todo {
	case "GETTASK":
		return worker.handleGETTASK(bag)
	case "REPORT":
		return worker.handleREPORT(bag)
	case "PING":
		return worker.handlePING(bag)
	default:
		return errors.New("bad request heartbeat todo")
	}

}

func (worker *Worker) handleGETTASK(bag *protocol.HeartBeatBag) error {
	var task = worker.ctx.TaskPool.GetTask()
	var pack *protocol.Package
	var err error
	if task == nil {
		pack, err = worker.makePingPack()
	} else {
		worker.Attach(task)
		pack, err = worker.makeRunTaskPack(task)
	}
	if err != nil {
		return err
	}
	return worker.sendPackage(pack)
}

func (worker *Worker) handleREPORT(bag *protocol.HeartBeatBag) error {
	var report = bag.Report
	worker.relTask.UpdateState(report)
	// has flag to stop work
	var pack *protocol.Package
	var err error
	var status = worker.relTask.Job.GetStatus()
	// when stopping task or task fail, stop other relative task
	if status == 0x03 || status == 0x06 {
		pack, err = worker.makeStopTaskPack()
	} else {
		pack, err = worker.makePingPack()
	}
	if err != nil {
		return err
	}
	worker.dealWithStatus(report)
	return worker.sendPackage(pack)
}

func (worker *Worker) dealWithStatus(state *job.TaskState) {
	switch state.Status {
	case "DONE":
		if worker.relTask.Job.TaskDone() {
			if worker.relTask.Job.CompareStatusAndSwap(0x05, 0x02) {
				worker.ctx.Store.UpdateJob(worker.relTask.Job)
			}
		}
		worker.Dettach()
	case "STOPPED":
		if worker.relTask.Job.DecRunning() == 0 {
			if worker.relTask.Job.CompareStatusAndSwap(0x04, 0x03) {
				worker.ctx.Store.UpdateJob(worker.relTask.Job)
			}
		}
		worker.Dettach()
	case "FAILED":
		if worker.relTask.Job.CompareStatusAndSwap(0x06, 0x02, 0x01) {
			worker.ctx.Store.UpdateJob(worker.relTask.Job)
		}
		worker.relTask.Job.DecRunning()
		worker.Dettach()
	default:
		var oldState = worker.relTask.GetState()
		if oldState.FrameAt < state.FrameAt {
			worker.relTask.Job.IncFinishFrames(state.FrameAt - oldState.FrameAt)
			worker.relTask.UpdateState(state)
		}
	}
}

func (worker *Worker) handlePING(bag *protocol.HeartBeatBag) error {
	worker.lastUpdate = time.Now().Unix()
	pack, err := worker.makePingPack()
	if err != nil {
		return err
	}
	return worker.sendPackage(pack)
}

func (worker *Worker) makeRunTaskPack(t *job.Task) (*protocol.Package, error) {
	var bag = new(protocol.HeartBeatBag)
	bag.Todo = "RUNTASK"
	bag.Task = t
	return worker.makePackWithBag(bag)
}

func (worker *Worker) makeStopTaskPack() (*protocol.Package, error) {
	var bag = new(protocol.HeartBeatBag)
	bag.Todo = "STOPTASK"
	return worker.makePackWithBag(bag)
}

func (worker *Worker) makePingPack() (*protocol.Package, error) {
	var bag = new(protocol.HeartBeatBag)
	bag.Todo = "PING"
	return worker.makePackWithBag(bag)
}

func (worker *Worker) makePackWithBag(bag *protocol.HeartBeatBag) (*protocol.Package, error) {
	payload, err := json.Marshal(bag)
	if err != nil {
		return nil, err
	}
	return protocol.NewPackageWithPayload(0x02, payload), nil
}

func (worker *Worker) sendPackage(pack *protocol.Package) error {
	message, err := pack.Marshal()
	if err != nil {
		return err
	}
	if _, err := worker.conn.Write(message); err != nil {
		return err
	}
	return nil
}

func (worker *Worker) receivePackage() (*protocol.Package, error) {
	var pack = new(protocol.Package)
	if err := pack.Unmarshal(worker.reader); err != nil {
		return nil, err
	}
	return pack, nil
}
