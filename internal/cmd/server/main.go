package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/jencisoll/TCP-Chat-Server-RAW/internal/protocol"
	"github.com/jencisoll/TCP-Chat-Server-RAW/internal/server"
)

func main() {
	//Iniciar el servidor en el puerto 9000
	if err := server.Start("9000"); err != nil {
		log.Fatal("Error al iniciar servidor: %v", err)
	}
	//3. el hilo principal va a leer el teclado y lo empaquetara
	tecladoScanner := bufio.NewScanner(os.Stdin)
	for tecladoScanner.Scan() {
		textoCrudo := tecladoScanner.Text()

		//armamos el sobre oficial usando nuestra estructura
		msg := protocol.Message{
			Type:    protocol.CmdMsg,
			From:    "Jerry",    //aqui luego pediremos el  nombre al inicio
			Room:    "#General", //van a la sala general
			Payload: textoCrudo, // el mensaje real
		}

		//Convertimos la estructura a bytes JSON
		jsonBytes, err := json.Marshal(msg)
		if err != nil {
			fmt.Println("Error al serializando servidor: %v", err)
			continue
		}

		//enviamos el JSON por el cable, agregando el salto de línea vital
		conn.Write(append(jsonBytes, '\n'))
	}
}
