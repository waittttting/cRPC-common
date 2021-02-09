package tcp

const (
	// 消息
	businessesMsg = 1
	kitMsg = 2
	// 消息类型
	cmdRegister = "register"
	// 数据帧头长度
	msgHeaderLength uint16 = 300
)

const kitMsgRegister = "kitMsgRegister"

type Message struct {
	Header *Header
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
	// 消息类型
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

func NewRegisterMessage(serverVersion string, serverName string) *Message {

	payload := make([]byte, 300)
	for i := range payload {
		payload[i] = 255
	}
	header := &Header{
		PayloadLen: uint16(len(payload)),
		MsgCode: kitMsg,
		Cmd: kitMsgRegister,
		ServerVersion: serverVersion,
		ServerName: serverName,
	}


	return &Message{
		Header: header,
		Payload: &payload,
	}
}
