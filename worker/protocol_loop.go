package worker

import (
	"errors"
	"log"
	"net"

	"bufio"

	"encoding/json"

	"time"

	"github.com/krufyliu/dkvgo/protocol"
)

// ProtocolLoop define
type ProtocolLoop struct {
	ctx    *DkvWorker
	reader *bufio.Reader
}

// IOLoop implements ProtocolIO interface
func (loop *ProtocolLoop) IOLoop(conn net.Conn) error {
	loop.ctx.connection = conn
	loop.reader = bufio.NewReader(conn)
	if err := loop.register(); err != nil {
		return err
	}
	return loop.runLoop()
}

func (loop *ProtocolLoop) register() error {
	var pack *protocol.Package
	pack = protocol.NewPackage(0x01)
	if err := loop.sendPackage(pack); err != nil {
		return err
	}
	answer, err := loop.receivePackage()
	if err != nil {
		return err
	}
	if answer.Directive != 0x01 {
		return errors.New("bad register response")
	}
	var sessionID = string(answer.Payload)
	loop.ctx.sessionID = sessionID
	loop.ctx.joined = true
	log.Printf("register successfully, session: %x\n", sessionID)
	return nil
}

func (loop *ProtocolLoop) runLoop() error {
	// var pack *protocol.Package
	// var err error
	// for {
	// 	pack, err = loop.receivePackage()
	// 	if err != nil {
	// 		break
	// 	}
	// 	err = loop.execDirective(pack)
	// 	if err != nil {
	// 		break
	// 	}
	// }
	// return err
	for {
		if err := loop.makeOnePullRequest(); err != nil {
			return err
		}
		time.Sleep(10 * time.Second)
	}

}

func (loop *ProtocolLoop) makeOnePullRequest() error {
	for {
		var reqpack = loop.newHeartBeatPack()
		if err := loop.sendPackage(reqpack); err != nil {
			if value, ok := err.(net.Error); ok && value.Timeout() {
				time.Sleep(2)
				continue
			} else {
				return err
			}
		}
		break
	}

	for {
		resPack, err := loop.receivePackage()
		if err != nil {
			if value, ok := err.(net.Error); ok && value.Timeout() {
				time.Sleep(2)
				continue
			} else {
				return err
			}
		}
		return loop.handleResponsePack(resPack)
	}
}

func (loop *ProtocolLoop) newHeartBeatPack() *protocol.Package {
	heartMessage, _ := json.Marshal(loop.newHeartBeatBag())
	var pack = protocol.NewPackageWithPayload(0x02, heartMessage)
	return pack
}

func (loop *ProtocolLoop) newHeartBeatBag() *protocol.HeartBeatBag {
	var heartb = new(protocol.HeartBeatBag)
	var ctx = loop.ctx.ctx
	if ctx == nil {
		heartb.Todo = "GETTASK"
		log.Printf("make heartbeat with todo %s\n", heartb.Todo)
	} else if ctx.state != nil {
		heartb.Todo = "REPORT"
		heartb.Report = ctx.state
		if ctx.state.Status != "RUNNING" {
			loop.ctx.clearCtx()
		}
		log.Printf("make heartbeat with todo %s: %s\n", heartb.Todo, heartb.Report)
	} else {
		heartb.Todo = "PING"
		log.Printf("make heartbeat with todo %s\n", heartb.Todo)
	}
	return heartb
}

func (loop *ProtocolLoop) handleResponsePack(pack *protocol.Package) error {
	if pack.WithPack != 1 && pack.Directive != 0x02 {
		return errors.New("bad heartbeat response")
	}
	var heartb = new(protocol.HeartBeatBag)
	if err := json.Unmarshal(pack.Payload, heartb); err != nil {
		return err
	}
	log.Printf("receive heartbeat with todo %s\n", heartb.Todo)
	switch heartb.Todo {
	case "RUNTASK":
		log.Printf("get task %s\n", heartb.Task)
		go loop.ctx.runTask(heartb.Task)
		loop.ctx.updatePing()
	case "STOPTASK":
		loop.ctx.forceStopTask()
		loop.ctx.updatePing()
	case "PING":
		loop.ctx.updatePing()
	default:
		log.Printf("Error: Unkown heartbeat %s\n", heartb.Todo)
		return errors.New("bad heartbeat todo")
	}
	return nil
}

// func (loop *ProtocolLoop) execDirective(pack *protocol.Package) error {
// 	switch pack.Directive {
// 	case protocol.TaskSubmit:
// 		return loop.handleTaskSubmit(pack)
// 	case protocol.TaskStop:
// 		return loop.handleTaskStop(pack)
// 	default:
// 		// should never reach here
// 		return errors.New("bad directive")
// 	}
// }

// func (loop *ProtocolLoop) handleTaskSubmit(pack *protocol.Package) error {
// 	var task = new(job.Task)
// 	if err := json.Unmarshal(pack.Payload, task); err != nil {
// 		return err
// 	}
// 	loop.ctx.setSessionState(TaskSubmitAccepted)
// 	loop.ctx.runTask(task)
// 	var answer = protocol.NewPackage(protocol.TaskSumbitAccept)
// 	return loop.sendPackage(answer)
// }

// func (loop *ProtocolLoop) handleTaskStop(pack *protocol.Package) error {
// 	loop.ctx.setSessionState(TaskStopAccepted)
// 	loop.ctx.stopTask()
// 	return nil
// }

func (loop *ProtocolLoop) sendPackage(pack *protocol.Package) error {
	message, err := pack.Marshal()
	if err != nil {
		return err
	}
	if _, err := loop.ctx.connection.Write(message); err != nil {
		return err
	}
	return nil
}

func (loop *ProtocolLoop) receivePackage() (*protocol.Package, error) {
	var pack = new(protocol.Package)
	if err := pack.Unmarshal(loop.reader); err != nil {
		return nil, err
	}
	return pack, nil
}
