package protocol

import "encoding/json"

type CmdType string

const (
	CmdAuth CmdType = "AUTH" // Para registrar el nickname al conectar
	CmdMsg          = "MSG"  //Para mensajes de chat normales
	CmdJoin         = "JOIN" //Para unirse a una sala especifica
)

type Message struct {
	Type      CmdType `json:"type"`                //¿Que acción es? (ej.MSG)
	From      string  `json:"from,omitempty"`      //¿Quien lo envía? (tu nickname)
	Room      string  `json:"room,omitempty"`      // ¿A que sala va? (ej general)
	Payload   string  `json:"payload,omitempty"`   //¿El contenido real del mensaje
	Timestamp int64   `json:"timestamp,omitempty"` //Con timestamp, el cliente puede reordenar y mostrarlos correctamente
}

func (m *Message) ToJSON() ([]byte, error) {
	return json.Marshal(m)
}
func NewMessage() *Message {}
