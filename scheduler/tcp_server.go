package scheduler

import (
	"log"
	"net"
	"runtime"
)

type TCPHandler interface {
	Handle(net.Conn)
}

func TCPServer(listener net.Listener, handler TCPHandler) {
	for {
		clientConn, err := listener.Accept()
		if err != nil {
			if nerr, ok := err.(net.Error); ok && nerr.Temporary() {
				log.Printf("NOTICE: temporary Accept() failure - %s\n", err)
				runtime.Gosched()
				continue
			}
			log.Printf("ERROR: listener.Accept() - %s\n", err)
			break
		}
		log.Printf("client from %s connected\n", clientConn.RemoteAddr().String())
		go handler.Handle(clientConn)
	}
}
