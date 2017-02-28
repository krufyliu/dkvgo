package worker

import "net"

type IOLoop struct {
	context *DkvWorker
	conn    *net.Conn
}

func (ioloop *IOLoop) IOLoop(conn *net.Conn) error {
	ioloop.conn = conn
	return nil
}

func (ioloop *IOLoop) Register() error {

}
