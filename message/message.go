package message

import (
	"bytes"
	"encoding/binary"
	"errors"
)

type Message struct {
	Type  int
	Data  []byte
	Reply chan Message // Message到达之后，可能需要根据该Channel返回确认消息
}

func NewMessage(messageType int) (*Message, error) {
	m := &Message{}
	if messageType != MessageTypeConfirmBlock &&
		messageType != MessageTypeConfirmTransaction &&
		messageType != MessageTypeSendBlock &&
		messageType != MessageTypeSendTransaction &&
		messageType != MessageTypePort {
		return m, errors.New("type error")
	}
	m.Type = messageType
	return m, nil
}

func (m *Message) MarshalBinary() []byte {
	var newB = make([]byte, 4)
	binary.LittleEndian.PutUint32(newB, uint32(m.Type))

	buf := &bytes.Buffer{}
	buf.Write(newB)
	buf.Write(m.Data)
	return buf.Bytes()
}

func (m *Message) UnmarshalBinary(data []byte) {
	buf := bytes.NewBuffer(data)
	m.Type = int(binary.LittleEndian.Uint32(buf.Next(4)))
	m.Data = buf.Next(len(data) - 1)
	return
}
