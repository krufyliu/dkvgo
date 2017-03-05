package scheduler

import (
	"net"
	"time"
	"encoding/json"
	"bufio"
	"errors"
	"crypto/md5"
	"github.com/krufyliu/dkvgo/job"
	"github.com/krufyliu/dkvgo/protocol"
)

type Worker struct {
	ctx        *DkvScheduler
	conn       net.Conn
	reader 	   *bufio.Reader
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
		worker.conn.Close()
		return
	}
	worker.runLoop()
}

func (worker *Worker) identity() error {
	var deadline = time.Now().Add(2*time.Second)
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

}

func (worker *Worker) makeOnePullRequest() error {
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
	var heartb = new(protocol.HeartBeatBag)
	err := json.Unmarshal(pack.Payload, heartb)
	if err != nil {
		return nil
	}
	switch heartb.Todo {
	case "GETTASK":
		return worker.handleGETTASK()
	case "REPORT":
		return worker.handleREPORT()
	case "PING":
		return worker.handlePING()
	default:
		return errors.New("bad request heartbeat todo")
	}

}

func (worker *Worker) handleGETTASK() error {
	return nil
}

func (worker *Worker) handleREPORT() error {
	return nil
}

func (worker *Worker) handlePING() error {
	return nil
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
