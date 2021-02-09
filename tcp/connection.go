package tcp

import (
	"bytes"
	"encoding/binary"
	"github.com/sirupsen/logrus"
	"github.com/waittttting/cRPC-common/assert"
	"io"
	"net"
	"time"
)

type Connection struct {
	socket net.Conn
}

func NewConnection(socket net.Conn) *Connection {
	return &Connection{
		socket: socket,
	}
}

func CreateSocket(host string) *Connection {

	socket, err := net.DialTimeout("tcp", host, 3 * time.Second)
	if err != nil {
		logrus.Fatalf("create socket error [%v]", err)
	}
	return &Connection{
		socket: socket,
	}
}

func (conn *Connection) Send(msg *Message) error {

	msgBuf := bytes.NewBuffer(make([]byte, 0, msgHeaderLength + msg.Header.PayloadLen))
	// header buf
	var headerBytes [msgHeaderLength]byte

	copy(headerBytes[:32],            msg.Header.Uid)
	copy(headerBytes[32:32 + 32],     msg.Header.Token)
	copy(headerBytes[64:64 + 15],     msg.Header.Ip)
	copy(headerBytes[79:79 + 1],      []uint8{msg.Header.Tag})
	copy(headerBytes[80:80 + 32],     msg.Header.MTest)
	copy(headerBytes[112:112 + 32],   msg.Header.SessionId)
	copy(headerBytes[144:144 + 1],    []uint8{msg.Header.MsgCode})
	copy(headerBytes[145:145 + 32],   msg.Header.MsgId)
	copy(headerBytes[177:177 + 32],   msg.Header.TraceId)
	copy(headerBytes[209:209 + 32],   msg.Header.ServerName)
	copy(headerBytes[241:241 + 32],   msg.Header.ServerVersion)
	copy(headerBytes[273:273 + 2],    []uint8{uint8(msg.Header.PayloadLen >> 8), uint8(msg.Header.PayloadLen)})
	copy(headerBytes[275:275 + 25],   msg.Header.Expend)

	err := binary.Write(msgBuf, binary.BigEndian, &headerBytes)
	if err != nil {
		logrus.Errorf("write packet header to buf error [%v]", err)
		return err
	}
	// payload buf
	if msg.Header.PayloadLen != 0 {
		err = binary.Write(msgBuf, binary.BigEndian, msg.Payload)
		if err != nil {
			logrus.Errorf("write packet payload to buf error [%v]", err)
			return err
		}
	}

	err = socketSend(conn.socket, int64(len(msgBuf.Bytes())), msgBuf.Bytes(), 0)
	if err != nil {
		return err
	}
	return nil
}

func socketSend(conn net.Conn, packetLen int64, buffer []byte, timeout time.Duration) error {

	assert.CheckParam(packetLen <= int64(len(buffer)) && timeout >= 0)
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

func (conn *Connection) Receive(timeout time.Duration) (*Message, error) {
	msg := new(Message)
	headerBuff := make([]byte, msgHeaderLength)

	err := socketReceive(conn.socket, int64(msgHeaderLength), headerBuff, timeout)
	if err != nil {
		logrus.Warningf("read header err, socket : [%s], error [%v]", conn.socket.RemoteAddr(), err)
		return nil, err
	}
	header := new(Header)
	header.ServerName =    string(bytes.Trim(headerBuff[209:209 + 32], "\x00"))
	header.ServerVersion = string(bytes.Trim(headerBuff[241:241 + 32], "\x00"))
	header.PayloadLen = uint16(headerBuff[273]) << 8 + uint16(headerBuff[274])

	payloadBuff := make([]byte, header.PayloadLen)
	err = socketReceive(conn.socket, int64(header.PayloadLen), payloadBuff, timeout)
	if err != nil {
		logrus.Warningf("read paylod err, socket : [%s], error [%v]", conn.socket.RemoteAddr(), err)
		return nil, err
	}
	msg.Header = header
	msg.Payload = &payloadBuff
	return msg, nil
}

func socketReceive(conn net.Conn, packetLen int64, Buffer []byte, timeout time.Duration)  (err error) {

	var receiveLen int64
	if timeout <= 0 {
		conn.SetReadDeadline(time.Time{})
	} else {
		conn.SetReadDeadline(time.Now().Add(timeout))
	}
	for receiveLen < packetLen {
		tempLen, err := conn.Read(Buffer[receiveLen:packetLen])
		if err == io.EOF {
			logrus.Warningf("peer socket[%s] exit", conn.RemoteAddr())
			return err
		} else if err != nil {
			logrus.Warningf("read socket[%s] connection error[%v]", conn.RemoteAddr(), err)
			return err
		} else {
			receiveLen += int64(tempLen)
		}
	}
	return err
}

