package scheduler

import (
	"net"
	"time"

	"crypto/md5"

	"github.com/krufyliu/dkvgo/protocol"
)

type ProtocolLoop struct {
	ctx        *DkvScheduler
	clientConn net.Conn
}

func (loop *ProtocolLoop) Handle(conn net.Conn) {
	loop.clientConn = conn
	if err := loop.identity(); err != nil {
		return
	}
	loop.ctx.Pool.Add(&Worker{conn: conn, remoteAddr: conn.RemoteAddr().String()})
}

func (loop *ProtocolLoop) identity() error {
	deadline = time.Now() + 2*time.Second
	loop.clientConn.SetReadDeadline(deadline)
	var pack = new(protocol.Package)
	if err := pack.Unmarshal(loop.clientConn); err != nil {
		return err
	}

	hash := md5.New()
	hash.Write([]byte(loop.clientConn.RemoteAddr().String()))
	hash.Write(time.Now().String())
	pack = protocol.NewPackageWithPayload(protocol.JoinAccept, hash.Sum(nil))
	message, err := pack.Marshal()
	if err != nil {
		return err
	}
	if _, err := loop.clientConn.Write(message); err != nil {
		return err
	}
	return nil
}
