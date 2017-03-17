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
	"github.com/krufyliu/dkvgo/scheduler/tracker"
)

type Worker struct {
	ctx        *DkvScheduler
	conn       net.Conn
	reader     *bufio.Reader
	remoteAddr string
	lastUpdate int64
	relTask    *job.Task
}

func (worker *Worker) RemoteAddr() string {
	if worker.conn == nil {
		return ""
	}
	return worker.conn.RemoteAddr().String()
}

func (worker *Worker) Attach(ts *job.Task) {
	worker.relTask = ts
	tracker.TraceTask(worker.relTask, worker.RemoteAddr())
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
	log.Printf("waiting for worker[%s] ident\n", worker.conn.RemoteAddr().String())
	var deadline = time.Now().Add(2 * time.Second)
	if err := worker.conn.SetReadDeadline(deadline); err != nil {
		return err
	}
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
	if err := worker.conn.SetReadDeadline(time.Time{}); err != nil {
		return err
	}
	log.Printf("worker[%s] ident successfully\n", worker.conn.RemoteAddr().String())
	return nil
}

func (worker *Worker) runLoop() {
	for {
		err := worker.makeOnePullResponse()
		if err != nil {
			if err == io.EOF {
				log.Printf("%s disconnected\n", worker.conn.RemoteAddr())
			} else {
				log.Printf("Error: %s\n", err)
			}
			worker.conn.Close()
			if worker.relTask != nil {
				log.Printf("push task %s back to task pool", worker.relTask)
				worker.ctx.TaskPool.PushFront(worker.relTask)
			}
			break
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
	log.Printf("%s receive heartbeat package with todo %s\n", worker.conn.RemoteAddr(), bag.Todo)
	switch bag.Todo {
	case "GETTASK":
		return worker.handleGETTASK(bag)
	case "REPORT":
		return worker.handleREPORT(bag)
	case "PING":
		return worker.handlePING(bag)
	default:
		return errors.New("bad request heartbeat todo " + bag.Todo)
	}
}

func (worker *Worker) handleGETTASK(bag *protocol.HeartBeatBag) error {
	var task = worker.ctx.TaskPool.GetTask()
	var pack *protocol.Package
	var err error
	if task == nil {
		log.Printf("%s get no task\n", worker.conn.RemoteAddr())
		pack, err = worker.makePingPack()
	} else {
		log.Printf("%s get task %s\n", worker.conn.RemoteAddr(), task)
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
	worker.dealWithStatus(report)
	// has flag to stop work
	var pack *protocol.Package
	var err error
	// when stopping task or task fail, stop other relative task
	if worker.relTask != nil {
		var status = worker.relTask.Job.GetStatus()
		if status == 0x03 || status == 0x04 || status == 0x06 {
			pack, err = worker.makeStopTaskPack()
		} else {
			pack, err = worker.makePingPack()
		}
	} else {
		pack, err = worker.makePingPack()
	}
	if err != nil {
		return err
	}
	return worker.sendPackage(pack)
}

func (worker *Worker) dealWithStatus(state *job.TaskState) {
	log.Printf("%s %s report: %s\n", worker.RemoteAddr(), worker.relTask, state)
	tracker.TraceTaskWithState(worker.relTask, worker.RemoteAddr(), state)
	if state.Status != "RUNNING" {
		worker.Dettach()
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
