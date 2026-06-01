package main

import (
	"log"

	"github.com/jencisoll/TCP-Chat-Server-RAW/internal/server"
)

func main() {

	//Iniciar el servidor en el puerto  9000
	if err := server.Start("9000"); err != nil {
		log.Fatal(err)
	}
}
