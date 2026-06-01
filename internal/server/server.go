package server

import (
	"fmt"
	"net"
)

// Start inicia el servidor TCP en el puerto indicado.
func Start(port string) error {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return fmt.Errorf("No se pudo iniciar el listener en puerto %s: %w", port, err)
	}
	defer listener.Close()
	fmt.Printf("Servidor TCP escuchando en puerto %s\n", port)
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Error al aceptar conexión: %v\n", err)
			continue
		}
		//cada conexión se maneja en su propia gorutine
		go handleConnection(conn)

	}
}

// handleconnection maneja una conexión individual
func handleConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Printf("Nueva conexión desde %s\n", conn.RemoteAddr())
	//aQUI VA la logica del chat (lectura/escritura)
}
