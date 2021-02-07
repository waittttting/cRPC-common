package tcp

import (
	"bytes"
	"encoding/binary"
	"github.com/sirupsen/logrus"
	"io"
	"net"
	"time"
)

type Connection struct {
	socket *net.Conn
}

func CreateSocket(host string) *Connection {

	socket, err := net.DialTimeout("tcp", host, 3 * time.Second)
	if err != nil {
		logrus.Fatalf("create socket error [%v]", err)
	}
	return &Connection{
		socket: &socket,
	}
}

func (con *Connection) Send(msg *Message) error {

	buf := bytes.NewBuffer(make([]byte, 0, msgHeaderLength + msg.header.payloadLen))
	err := binary.Write(buf, binary.BigEndian, &msg.header)
	if err != nil {
		logrus.Errorf("write packet header to buf error [%v]", err)
		return err
	}

	if msg.header.payloadLen != 0 {
		err = binary.Write(buf, binary.BigEndian, msg.payload)
		if err != nil {
			logrus.Errorf("write packet payload to buf error [%v]", err)
			return err
		}
	}
	err = SocketSend(*con.socket, int64(len(buf.Bytes())), buf.Bytes(), 0)
	if err != nil {
		return err
	}
	return nil
}

func SocketSend(conn net.Conn, packetLen int64, buffer []byte, timeout int64) error {

	hasWrittenLen := int64(0)
	for hasWrittenLen < packetLen {
		n, err := conn.Write(buffer[hasWrittenLen : packetLen])
		if err != nil {
			logrus.Errorf("conn write buffer error [%v]", err)
			return err
		}
		if n == 0 {
			// 一般不会出现 n == 0 的情况，如果出现，则可认定为不正常状态，打印告警日志，且返回 EOF 即可；
			//  Reference: http://stackoverflow.com/questions/2176443/is-a-return-value-of-0-from-write2-in-c-an-error
			logrus.Warningf("write socket[%s] zero bytes!", conn.RemoteAddr())
			return io.EOF
		} else {
			hasWrittenLen = hasWrittenLen + int64(n)
		}
	}
	return nil
}

func (con *Connection) Receive() {

}

