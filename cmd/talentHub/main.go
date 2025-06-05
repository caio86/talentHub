package main

import (
	"fmt"
	"log"

	"github.com/caio86/talentHub/http"
	"github.com/caio86/talentHub/postgres"
)

func main() {
	port := 8080
	svr := http.NewServer()

	// Setting Listener port
	svr.Addr = fmt.Sprintf(":%d", port)

	// Setting services
	candidatosService := postgres.NewCandidatoService()

	svr.CandidatoService = candidatosService

	log.Printf("listening: port=%d", port)
	if err := svr.Open(); err != nil {
		panic(err)
	}
}
