package tcp

const (
	// 消息
	businessesMsg = 1
	kitMsg = 2
	// 数据帧头长度
	msgHeaderLength uint16 = 4
)

const kitMsgRegister = "kitMsgRegister"

type Message struct {
	header *Header
	payload *[]byte
}

type Header struct {
	// 用户id
	uid [64]byte
	// token
	token [64]byte
	// IP
	ip [32]byte
	// 流量标记(用于灰度等)
	tag uint8
	// 全链路压测
	mTest [32]byte
	// 会话ID
	sessionId [32]byte
	// 消息类型
	msgCode uint8
	// 消息ID
	msgId [32]byte
	// 命令号
	cmd [64]byte
	// traceId 链路追踪
	traceId [64]byte
	// 服务名称
	serverName [32]byte
	// 服务版本
	serverVersion [32]byte
	// body 长度，在读数据的时候有用
	payloadLen uint16
}

func NewRegisterMessage(serverVersion string, serverName string) *Message {

	header := &Header{
		payloadLen: 0,
		msgCode: kitMsg,
		//uid: []byte(serverVersion),
	}

	return &Message{
		header: header,
		payload: nil,
	}
}