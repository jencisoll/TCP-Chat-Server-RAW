package server

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"sync"

	"github.com/jencisoll/TCP-Chat-Server-RAW/internal/protocol"
)

//Server mantiene  el estado global del chat

type Server struct {
	mu      sync.Mutex
	clients map[net.Conn]string //conexión ->nombre
}

//Start crea el servidor y comienza a escuchar en el puerto indicado

func Start(port string) error {
	s := &Server{
		clients: make(map[net.Conn]string),
	}
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return fmt.Errorf("Error al iniciar servidor: %w", err)
	}
	defer listener.Close()

	fmt.Printf("Servidor TCP escuchando en puerto %s\n", port)
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Error aceptando conexión: %v\n", err)
			continue
		}
		go s.handleConnection(conn)

	}
}

// HandleConnection administra la vida de un cliente.
func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	//asignar un nombre temporal basado en la dirección remota
	clientName := conn.RemoteAddr().String()

	s.mu.Lock()
	s.clients[conn] = clientName
	s.mu.Unlock()

	//Anunciar entrada
	s.broadcast(conn, fmt.Sprintf("[%s] se ha unido al chat", clientName))
	fmt.Printf("[NUEVA CONEXIÓN] %s conectado como %s\n", conn.RemoteAddr(), clientName)

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		textoCrudo := scanner.Text()

		//1 preparamos el molde vacio donde volcaremos los datos
		var msg protocol.Message

		//2. Intentamos abrir el paquete JSON y volcarlo en el  molde
		err := json.Unmarshal([]byte(textoCrudo), &msg)
		if err != nil {
			fmt.Printf("[ADVERTENCIA] El cliente %s enció datos que no son JSON: %s\n", clientName, textoCrudo)
			continue //ignoramos este mensaje y esperamos el siguiente
		}
		switch msg.Type {
		case protocol.CmdMsg:
			fmt.Printf("[Sala: %s] %s dice: %s\n", msg.Room, msg.From, msg.Payload)
		case protocol.CmdJoin:
			fmt.Printf("[Sala: %s] %s acaba de entrar a la sala: %s\n", msg.From, msg.Room)
		default:
			fmt.Printf(" comando desconocido de %s : %s\n", msg.From, msg.Type)
		}
	}
	//Si sale del bucle el cliente se desconecto
	s.mu.Lock()
	delete(s.clients, conn)
	s.mu.Unlock()
	s.broadcast(conn, fmt.Sprintf("[%s] ha salido del chat", clientName))
	fmt.Printf("[DESCONEXIÓN] %s (%s) se fue\n", clientName, conn.RemoteAddr())
}

// broadcast envía un mensaje a todos los clientes excepto al remitente.
func (s *Server) broadcast(sender net.Conn, message string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for conn := range s.clients {
		if conn != sender {
			// Escribir mensaje + salto de línea
			fmt.Fprintln(conn, message)
		}
	}
}
