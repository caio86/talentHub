package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/caio86/talentHub/http"
	"github.com/caio86/talentHub/postgres"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	port := 8080
	svr := http.NewServer()

	// Setting Listener port
	svr.Addr = fmt.Sprintf(":%d", port)

	// Connecting to DB
	dbConfig, err := loadDBVars()
	if err != nil {
		log.Fatal(err)
	}
	db := postgres.NewDB(dbConfig)
	if err := db.Connect(); err != nil {
		log.Fatal(err)
	}

	// Setting services
	candidatosService := postgres.NewCandidatoService(db)
	vagaService := postgres.NewVagaService(db)
	applicationService := postgres.NewApplicationService(db)

	svr.CandidatoService = candidatosService
	svr.VagaService = vagaService
	svr.ApplicationService = applicationService

	log.Printf("listening: port=%d", port)
	if err := svr.Open(); err != nil {
		panic(err)
	}
}

func loadPassword() (string, error) {
	password, ok := os.LookupEnv("DB_PASSWORD")
	if ok {
		return password, nil
	}

	passwordFile, ok := os.LookupEnv("DB_PASSWORD_FILE")
	if !ok {
		return "", fmt.Errorf("DB_PASSWORD or DB_PASSWORD_FILE env vars not set")
	}

	data, err := os.ReadFile(passwordFile)
	if err != nil {
		return "", fmt.Errorf("failed to read password file: %w", err)
	}

	return strings.TrimSpace(string(data)), nil
}

func loadDBVars() (*postgres.DBConfig, error) {
	username, ok := os.LookupEnv("DB_USER")
	if !ok {
		return nil, fmt.Errorf("DB_USER env var not set")
	}

	password, err := loadPassword()
	if err != nil {
		return nil, fmt.Errorf("loading password: %w", err)
	}

	host, ok := os.LookupEnv("DB_HOST")
	if !ok {
		return nil, fmt.Errorf("DB_HOST env var not set")
	}

	portStr, ok := os.LookupEnv("DB_PORT")
	if !ok {
		return nil, fmt.Errorf("DB_PORT env var not set")
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("failed to convert port to int: %w", err)
	}

	dbname, ok := os.LookupEnv("DB_NAME")
	if !ok {
		return nil, fmt.Errorf("DB_NAME env var not set")
	}

	sslmode, ok := os.LookupEnv("DB_SSLMODE")
	if !ok {
		return nil, fmt.Errorf("DB_SSLMODE env var not set")
	}

	config := &postgres.DBConfig{
		User:     username,
		Password: password,
		Host:     host,
		Port:     uint16(port),
		DBName:   dbname,
		SSLMode:  sslmode,
	}

	err = config.Validate()
	if err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return config, nil
}
