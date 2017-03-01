package worker

import (
	"errors"
	"net"

	"io"

	"bufio"

	"encoding/json"

	"github.com/krufyliu/dkvgo/protocol"
	"github.com/krufyliu/dkvgo/task"
)

// ProtocolLoop define
type ProtocolLoop struct {
	context *DkvWorker
	r       io.Reader
}

// IOLoop implements ProtocolIO interface
func (loop *ProtocolLoop) IOLoop(conn *net.TCPConn) error {
	loop.context.connection = conn
	return nil
}

func (loop *ProtocolLoop) reader() io.Reader {
	if loop.r == nil {
		loop.r = bufio.NewReader(loop.context.connection)
	}
	return loop.r
}

func (loop *ProtocolLoop) register() error {
	var pack *protocol.Package
	if loop.context.sessionID == "" {
		pack = protocol.NewPackage(protocol.Join)
	} else {
		pack = protocol.NewPackageWithPayload(protocol.Join, []byte(loop.context.sessionID))
	}
	if err := loop.sendPackage(pack); err != nil {
		return err
	}
	answer, err := loop.receivePackage()
	if err != nil {
		return err
	}
	if answer.Directive != protocol.JoinAccept {
		return errors.New("join failed")
	}
	var sessionID = string(answer.Payload)
	loop.context.setSessionID(sessionID)
	return nil
}

func (loop *ProtocolLoop) runLoop() error {
	var pack *protocol.Package
	var err error
	for {
		pack, err = loop.receivePackage()
		if err != nil {
			break
		}
		err = loop.execDirective(pack)
		if err != nil {
			break
		}
	}
	return err
}

func (loop *ProtocolLoop) execDirective(pack *protocol.Package) error {
	switch pack.Directive {
	case protocol.TaskSubmit:
		return loop.handleTaskSubmit(pack)
	case protocol.TaskStop:
		return loop.handleTaskStop(pack)
	default:
		// should never reach here
		return errors.New("bad directive")
	}
}

func (loop *ProtocolLoop) handleTaskSubmit(pack *protocol.Package) error {
	var taskSeg task.TaskSegment
	if err := json.Unmarshal(pack.Payload, &taskSeg); err != nil {
		return err
	}
	var answer = protocol.NewPackage(protocol.TaskSumbitAccept)
	return loop.sendPackage(answer)
}

func (loop *ProtocolLoop) handleTaskStop(pack *protocol.Package) error {
	return nil
}

func (loop *ProtocolLoop) sendPackage(pack *protocol.Package) error {
	message, err := pack.Marshal()
	if err != nil {
		return err
	}
	if _, err := loop.context.connection.Write(message); err != nil {
		return err
	}
	return nil
}

func (loop *ProtocolLoop) receivePackage() (*protocol.Package, error) {
	var pack = new(protocol.Package)
	if err := pack.Unmarshal(loop.reader()); err != nil {
		return nil, err
	}
	return pack, nil
}
