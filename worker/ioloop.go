package worker

import (
	"bufio"
	"encoding/json"
	"github.com/krufyliu/dkvgo"
	"github.com/krufyliu/dkvgo/protocol"
	"log"
	"net"
)

type IOLoop struct {
	context *DkvWorker
	reader  *bufio.Reader
}

func (ioloop *IOLoop) HandleRequest() error {
	var req *protocol.DkvRequst
	reader := ioloop.bufferReader()
	for {
		req = &protocol.DkvRequst{}
		if err := req.PullHeader(reader); err != nil {
			return err
		}
	}
}

func (ioloop *IOLoop) IOLoop(conn *net.Conn) error {
	if err := ioloop.Register(); err != nil {
		return err
	}
	return nil
}

func (ioloop *IOLoop) Register() error {
	var req *protocol.DkvRequst
	if ioloop.context.sessionID != "" {
		req = protocol.NewDkvRequest("REGISTER", []string{ioloop.context.sessionID})
	} else {
		req = protocol.NewDkvRequest("REGISTER", []string{})
	}
	err := ioloop.seedRequest(req)
	if err != nil {
		return nil
	}
	res, err := ioloop.expectResponse()
	if err != nil {
		return err
	}
	if res.Version != "V1.0" {
		return ErrVersionNoMatch
	}
	if res.Status != "OK" {
		return ErrResponseStatus
	}
	if err := res.PullContent(ioloop.bufferReader()); err != nil {
		return err
	}

	return nil
}

func (ioloop *IOLoop) RunTask(req *protocol.DkvRequst) {
	var taskSegment dkvgo.TaskSegment
	if err := json.Unmarshal(req.Payload, &taskSegment); err != nil {
		log.Fatalln(err)
	}
	ioloop.seedResponse(protocol.NewDkvResponse("OK"))
	cmdgen := dkvgo.NewCmdGeneratorFromTaskSegment(&taskSegment, 8, "/usr/local/visiondk/bin", "/usr/local/visiondk/setting")
	ioloop.context.taskSeg = &taskSegment
	ioloop.context.cmd = cmdgen.GetCmd()
}

func (ioloop *IOLoop) bufferReader() *bufio.Reader {
	if ioloop.reader != nil {
		return ioloop.reader
	}
	return bufio.NewReader(ioloop.context.connection)
}

func (ioloop *IOLoop) expectResponse() (*protocol.DkvResponse, error) {
	var res = &protocol.DkvResponse{}
	if err := res.PullHeader(ioloop.bufferReader()); err != nil {
		return nil, err
	}
	return res, nil
}

func (ioloop *IOLoop) acceptRequest() (*protocol.DkvRequst, error) {
	var req = &protocol.DkvRequst{}
	if err := req.PullHeader(ioloop.bufferReader()); err != nil {
		return nil, err
	}
	return req, nil
}

func (ioloop *IOLoop) seedResponse(res *protocol.DkvResponse) error {
	message, err := res.Dumps()
	if err != nil {
		return err
	}
	_, err = ioloop.context.connection.Write(message)
	return err
}

func (ioloop *IOLoop) seedRequest(req *protocol.DkvRequst) error {
	message, err := req.Dumps()
	if err != nil {
		return err
	}
	_, err = ioloop.context.connection.Write(message)
	return err
}
