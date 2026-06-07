package server

import (
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

	decoder := protocol.NewDecoder(conn)
	for {
		msg, err := decoder.Decode()
		if err != nil {
			break //Cliente se deconectó o error de lectura
		}
		switch msg.Type {
		case protocol.CmdMsg:
			fmt.Printf("[%s] %s:  %s\n", msg.Room, msg.From, msg.Payload)
			s.broadcast(conn, fmt.Sprintf("[%s]: %s", msg.From, msg.Payload))
		case protocol.CmdJoin:
			fmt.Printf("[%s] entró a %s\n", msg.From, msg.Room)
		default:
			fmt.Printf("comando desconocido:  %s\n", msg.Type)

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
			encoder := protocol.NewEncoder(conn)
			encoder.Encode(protocol.NewMessage(protocol.CmdMsg, "server", "#general", message))
		}
	}
}
