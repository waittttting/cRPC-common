package tcp

const (

	// 业务消息
	msgTypeBusinesses = 1

	// 框架消息
	msgTypekit = 2

	// client 注册到 control center
	msgKitRegisterPing = "msgKitRegisterPing"
	// control center 回复 client 的注册消息
	msgKitRegisterPong = "msgKitRegisterPong"

	// 数据帧头长度
	msgHeaderLength uint16 = 320
)

type Message struct {
	Header  *Header
	Payload *[]byte
}

type Header struct {
	// 用户id
	Uid string
	// token
	Token string
	// IP
	Ip string
	// 流量标记(用于灰度等)
	Tag uint8
	// 全链路压测
	MTest string
	// 会话ID
	SessionId string
	// 消息类型 (框架消息/业务消息)
	MsgCode uint8
	// 消息ID
	MsgId string
	// 命令号
	Cmd string
	// traceId 链路追踪
	TraceId string
	// 服务名称
	ServerName string
	// 服务版本
	ServerVersion string
	// 扩展字段
	Expend string
	// body 长度，在读数据的时候有用
	PayloadLen uint16
}

func MsgRegisterPing(serverVersion string, serverName string) *Message {

	header := &Header{
		PayloadLen:    0,
		MsgCode:       msgTypekit,
		Cmd:           msgKitRegisterPing,
		ServerVersion: serverVersion,
		ServerName:    serverName,
	}

	return &Message{
		Header:  header,
		Payload: nil,
	}
}

func MsgRegisterPong() *Message {

	header := &Header{
		PayloadLen: 0,
		MsgCode:    msgTypekit,
		Cmd:        msgKitRegisterPong,
	}
	return &Message{
		Header:  header,
		Payload: nil,
	}
}
