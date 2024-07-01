package ws

import (
	"github.com/gorilla/websocket"
	"github.com/nullexp/finman-gateway-service/pkg/infrastructure/http/protocol/model"
	"github.com/nullexp/finman-gateway-service/pkg/infrastructure/misc"
)

type DuplexConnection struct {
	conn   *websocket.Conn
	caller misc.Caller
}

func NewDuplexConnection(conn *websocket.Conn, caller misc.Caller) *DuplexConnection {
	return &DuplexConnection{conn: conn, caller: caller}
}

func (d DuplexConnection) Publish(topic string, message any) error {
	dto := NewJsonDtoMessage(topic, message)
	return d.conn.WriteJSON(dto)
}

func (d DuplexConnection) GetCaller() (misc.Caller, bool) {
	return d.caller, d.caller != nil
}

func (d DuplexConnection) ReadMessage() (messageType int, p []byte, err error) {
	return d.conn.ReadMessage()
}

func (d DuplexConnection) MustGetCaller() misc.Caller {
	if d.caller == nil {
		panic("caller not defined")
	}
	return d.caller
}

func (d DuplexConnection) SendError(code, message string) error {
	return d.Publish(ErrorTopic, model.DuplexError{Message: message, Code: code})
}

func (d DuplexConnection) Close() error {
	return d.conn.Close()
}