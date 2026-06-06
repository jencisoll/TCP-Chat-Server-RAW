package protocol

import (
	"encoding/json"
	"time"
)

type CmdType string

const (
	CmdAuth CmdType = "AUTH" // Para registrar el nickname al conectar
	CmdMsg  CmdType = "MSG"  //Para mensajes de chat normales
	CmdJoin CmdType = "JOIN" //Para unirse a una sala específica
)

type Message struct {
	Type      CmdType `json:"type"`                //¿Qué acción es? (ej. MSG)
	From      string  `json:"from,omitempty"`      //¿Quien lo envía? (tu nickname)
	Room      string  `json:"room,omitempty"`      // ¿A que sala va? (ej general)
	Payload   string  `json:"payload,omitempty"`   //¿El contenido real del mensaje?
	Timestamp int64   `json:"timestamp,omitempty"` //Con timestamp, el cliente puede reordenar y mostrarlos correctamente
}

func (m *Message) ToJSON() ([]byte, error) {
	return json.Marshal(m)
}
func NewMessage(cmdType CmdType, from, room, payload string) *Message {
	return &Message{
		Type:      cmdType,
		From:      from,
		Room:      room,
		Payload:   payload,
		Timestamp: time.Now().UnixMilli(),
	}
}

func ParseJSON(data []byte) (*Message, error) {
	var msg Message
	err := json.Unmarshal(data, &msg)
	if err != nil {
		return nil, err
	}
	return &msg, nil
}
