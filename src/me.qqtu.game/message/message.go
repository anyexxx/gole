package message

import (
	"code.google.com/p/goprotobuf/proto"
	"encoding/json"
	"errors"
	"me.qqtu.game/logger"
	"strings"
)

const (
	DATAFORMAT_JSON = 1
	DATAFORMAT_PB   = 2
)

//读取数据创建消息体
func ReadMessage(data []byte, format int) Message {
	var msgId uint16

	var idBytes []byte = data[0:2]
	var bodyBytes []byte = data[2:]

	msgId = 0
	msgId += uint16(idBytes[0])
	msgId = msgId << 8
	msgId += uint16(idBytes[1])

	switch format {
	case DATAFORMAT_JSON:
		return &JsonMessage{messageId: msgId, jsonStr: string(bodyBytes)}
	case DATAFORMAT_PB:
		return &PBMessage{messageId: msgId, bytes: bodyBytes}
	}
	return nil
}

type Message interface {
	MessageId() uint16
	ProtocolBytes() []byte
	Encode(interface{}) error //序列化对象，传入需要序列化对象指针
	Decode(interface{}) error //反序列化对象，传入需要保存结果的对象指针
}

//protocol buffer数据格式-----------------------------------------------
type PBMessage struct {
	messageId uint16
	bytes     []byte
}

func (msg *PBMessage) Encode(val interface{}) error {
	if data, ok := val.(proto.Message); ok {
		if res, err := proto.Marshal(data); err == nil {
			msg.bytes = res
			return err
		} else {
			msg.bytes = []byte{0}
			logger.GetLogger().Error(err.Error(), nil)
			return err
		}
	} else {
		err := errors.New("PBMessage.Encode error:val is not proto.Message")
		logger.GetLogger().Error(err.Error(), nil)
		return err
	}
	return nil
}

func (msg *PBMessage) Decode(ret interface{}) error {
	if data, ok := ret.(proto.Message); ok {
		return proto.Unmarshal(msg.bytes, data)
	} else {
		err := errors.New("PBMessage.Decode error:val is not proto.Message")
		logger.GetLogger().Error(err.Error(), nil)
		return err
	}
	return nil
}

func (msg *PBMessage) MessageId() uint16 {
	return msg.messageId
}

func (msg *PBMessage) ProtocolBytes() []byte {
	ret := make([]byte, 2)
	ret[0] = byte(msg.MessageId() >> 8)
	ret[1] = byte(msg.MessageId())
	ret = append(ret, msg.bytes...)
	return ret
}

func NewPBMessage(id uint16, data interface{}) *PBMessage {
	msg := &PBMessage{messageId: id}
	msg.Encode(data)
	return msg
}

//JSON格式消息-----------------------------------------------------------------
type JsonMessage struct {
	messageId uint16
	jsonStr   string
}

func (msg *JsonMessage) Encode(val interface{}) error {
	if res, err := json.Marshal(val); err == nil {
		msg.jsonStr = string(res)
		return nil
	} else {
		msg.jsonStr = ""
		logger.GetLogger().Error(err.Error(), nil)
		return err
	}
	return nil
}

func (msg *JsonMessage) Decode(ret interface{}) error {
	if strings.TrimSpace(msg.jsonStr) != "" {
		if err := json.Unmarshal([]byte(msg.jsonStr), ret); err != nil {
			logger.GetLogger().Error("JsonMessage Decode Error!", logger.Extras{"json": msg.jsonStr, "err": err})
			return err
		}
	}
	return nil
}

func (msg *JsonMessage) MessageId() uint16 {
	return msg.messageId
}

func (msg *JsonMessage) ProtocolBytes() []byte {
	ret := make([]byte, 2)
	ret[0] = byte(msg.MessageId() >> 8)
	ret[1] = byte(msg.MessageId())
	ret = append(ret, []byte(msg.jsonStr)...)
	return ret
}

func NewJsonMessage(id uint16, data interface{}) *JsonMessage {
	msg := &JsonMessage{messageId: id}
	msg.Encode(data)
	return msg
}
