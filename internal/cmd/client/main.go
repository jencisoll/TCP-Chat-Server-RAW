package main

import (
	"bufio"
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
	encoder := protocol.NewEncoder(conn)
	fmt.Println("🔌 Conectado. Escribe un mensaje:")
	// Antes del loop del teclado, lanza esto:
	go func() {
		decoder := protocol.NewDecoder(conn)
		for {
			msg, err := decoder.Decode()
			if err != nil {
				fmt.Println("Desconectado del servidor")
				os.Exit(0)
			}
			fmt.Printf("[%s] %s: %s\n", msg.Room, msg.From, msg.Payload)
		}
	}()
	tecladoScanner := bufio.NewScanner(os.Stdin)
	for tecladoScanner.Scan() {
		textoUsuario := tecladoScanner.Text()

		// 2. Lo convertimos a JSON
		encoder.Encode(protocol.NewMessage(protocol.CmdMsg, "Jerry", "#general", textoUsuario))

	}
}
