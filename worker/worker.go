package worker

import (
	"net"
)

// Worker define worker struct
type Worker struct {
	options    Options
	connection *net.TCPConn
	retry      int
	waitTime   int
}

func (w *Worker) connect() error {
	conn, err := net.DialTCP("tcp", nil, w.options.schedulerAddr)
	if err != nil {
		return err
	}
	w.connection = conn
	return nil
}

func (w *Worker) register() error {
	return nil
}
