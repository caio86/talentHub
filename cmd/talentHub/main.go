package main

import (
	"context"
	"fmt"
	"log"

	"github.com/caio86/talentHub/http"
	"github.com/caio86/talentHub/postgres"
	"github.com/jackc/pgx/v5"
)

func main() {
	port := 8080
	svr := http.NewServer()

	// Setting Listener port
	svr.Addr = fmt.Sprintf(":%d", port)

	// Connecting to DB
	host := "localhost"
	dbPort := "5432"
	user := "postgres"
	dbname := "talentHub"
	password := "password"

	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		host,
		dbPort,
		user,
		dbname,
		password,
	)

	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}

	// Setting services
	candidatosService := postgres.NewCandidatoService()
	vagaService := postgres.NewVagaService(conn)

	svr.CandidatoService = candidatosService
	svr.VagaService = vagaService

	log.Printf("listening: port=%d", port)
	if err := svr.Open(); err != nil {
		panic(err)
	}
}
