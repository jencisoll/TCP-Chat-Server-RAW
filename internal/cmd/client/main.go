package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"

	// OJO: Ajusta esta ruta a tu proyecto exacto
	"github.com/jencisoll/TCP-Chat-Server-RAW/internal/protocol"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:9000")
	if err != nil {
		log.Fatalf("No se pudo conectar: %v", err)
	}
	defer conn.Close()

	fmt.Println("🔌 Conectado. Escribe un mensaje:")

	tecladoScanner := bufio.NewScanner(os.Stdin)
	for tecladoScanner.Scan() {
		textoUsuario := tecladoScanner.Text()

		// 1. Armamos el sobre oficial
		msg := protocol.Message{
			Type:    protocol.CmdMsg,
			From:    "Jerry", // Tu nombre
			Room:    "#general",
			Payload: textoUsuario,
		}

		// 2. Lo convertimos a JSON
		jsonBytes, err := json.Marshal(msg)
		if err != nil {
			fmt.Println("Error empaquetando:", err)
			continue
		}

		// 3. Lo enviamos por el cable con el salto de línea vital (\n)
		conn.Write(append(jsonBytes, '\n'))
	}
}
