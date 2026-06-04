package protocol

import "encoding/json"

const (
	CmdAuth = "AUTH" // Para registrar el nickname al conectar
	CmdMsg  = "MSG"  //Para mensajes de chat normales
	CmdJoin = "JOIN" //Para unirse a una sala especifica
)

type Message struct {
	Type    string `json:"type"`    //¿Que acción es? (ej.MSG)
	From    string `json:"from"`    //¿Quien lo envía? (tu nickname)
	Room    string `json:"room"`    // ¿A que sala va? (ej general)
	Payload string `json:"payload"` //¿El contenido real del mensaje
}

func (m *Message) ToJson() ([]byte, error) {
	return json.Marshal(m)
}
