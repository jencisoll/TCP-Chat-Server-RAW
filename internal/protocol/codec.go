package protocol

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"io"

	"sync"
)

const MaxMessageSize = 64 * 1024

type Encoder struct {
	w  io.Writer  //La conexión TCP cruda.
	mu sync.Mutex //El candado para evitar choques si mandamos 2 cosas a la vez.
}

// Decoder se encarga de leer e interpretar datos desde un flujo de entrada.
type Decoder struct {
	r io.Reader //La fuente de origen (ej. el socket TCP)
}

// NewEncoder es el constructor que inicializa nuestro escritor.
func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{w: w}
}

// NewDecoder es el constructor que inicializa nuestro lector.
func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{r: r}
}

// Encode es el método que prepara y envía el mensaje por la red
func (c *Encoder) Encode(msg Message) error {
	//paso 1. Convertir mensaje a JSON
	//json.Marsal toma tu estructura amigable y la trutra hasta convertirla en un punto de bytes puros ([]byte)

	jsonData, err := json.Marshal(msg)
	if err != nil {
		//si la estructura tiene datos incompatibles, fallara aquí
		return err
	}
	//-------PASO 2 -------
	//Usamos len(jsonData) porque en redes nos importa el peso en bytes, No la cantidad de letra.
	if len(jsonData) > MaxMessageSize {
		return errors.New("mensaje exede el limite")
	}
	//Cerramos
	c.mu.Lock()

	//Inmediatamente aperturamos
	defer c.mu.Unlock()
	// ---------PASOS 4 Y5: Crear el paquete total y meter el tamaño----
	//len() devuelve un 'int', pero BigEndian exige estrictamente un 'uint32'. hacemos la conversion explisita (casting)

	tamañoMensaje := uint32(len(jsonData))

	//Enb lugar de hacer una cajita de 4 y enviar el JSON aparte, creamos el "paquete gigante" continuo en memoria
	//Tamaño = 4 bytes del encabezado + el peso del JSON
	paqueteFinal := make([]byte, 4+tamañoMensaje)

	//Inyectamos el número en formato Big-Endian
	//OJO a la sintaxis paqueteFinal[:4] -> le decimos a G:
	//"Toma este arreglo. Pero solo usa desde el principio hasta el indice 3".
	binary.BigEndian.PutUint32(paqueteFinal[:4], tamañoMensaje)

	//------PASO 6: Copiar el JSON al resto del paquete -------
	//La función nativa çopy(destino, origen)' es la herramientea más rápida
	//y segura en Go para mover bloques de bytes de un lado a otro en la memoria RAM
	//PaqueteFinal[4:] significa: "El destino empieza a apartir deel indice 4 hacia adelante"
	//Asi evitamos sobreescribir nuestro hermoso encabezado numérico.

	copy(paqueteFinal[4:], jsonData)

	//------PASO 7: El envio  a la red  (Llamada al sistema) ----
	//¡El gran momento! hacemos un único envio por la tarjeta de red.
	//y un error por si el cable virtula de cortó a la mitad del proceso

	_, err = c.w.Write(paqueteFinal)
	if err != nil {
		return errors.New("Fallo fatal al escribnir en el socket TCP: " + err.Error())

	}
	return nil
}

func (d *Decoder) Decode() (*Message, error) {

	//Paso 1: leer el encabezado fijo de 4  bytes---
	header := make([]byte, 4)

	//io.ReadFull garantiza que se lean los 4 bytes  completos antes de continuar.
	//Si el cliente se desconecta limpiamente, devolvera el error io.EOF.
	_, err := io.ReadFull(d.r, header)
	if err != nil {
		return nil, err
	}
	// --- PASO 2: Traducir los 4 bytes a la longitud real del mensaje ---
	// Pasamos de bytes en formato Network Byte Order (Big-Endian) a un número uint32.
	msgSize := binary.BigEndian.Uint32(header)

	// --- PASO 3: Control de seguridad ---
	// Si el tamaño extraído es sospechosamente gigante, rechazamos el procesamiento.
	if msgSize > MaxMessageSize {
		return nil, errors.New("mensaje entrante demasiado grande o corrupto")
	}

	// --- PASO 4: Leer el cuerpo del mensaje (El JSON crudo) ---
	// Creamos un buffer a la medida exacta del tamaño que leímos en el encabezado.
	payload := make([]byte, msgSize)

	// Leemos del flujo de red la cantidad exacta de bytes que corresponden al JSON.
	_, err = io.ReadFull(d.r, payload)
	if err != nil {
		return nil, err
	}

	// --- PASO 5: Reconstruir la estructura original ---
	var msg Message
	// json.Unmarshal toma los bytes del JSON y rellena los campos de la variable 'msg'.
	err = json.Unmarshal(payload, &msg)
	if err != nil {
		return nil, errors.New("error al deserializar el JSON: " + err.Error())
	}

	// Todo fue un éxito, retornamos la dirección del mensaje estructurado.
	return &msg, nil
}
