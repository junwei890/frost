package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/junwei890/rumbling/internal/database"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

type apiConfig struct {
	db *database.Queries
}

func main() {
	config := apiConfig{}

	err := godotenv.Load()
	if err != nil {
		log.Println("no environment variables loaded from .env file")
	}

	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		log.Fatal("no database url provided")
	}

	db, err := sql.Open("libsql", dbUrl)
	if err != nil {
		log.Fatal("no connection to database")
	}

	dbQueries := database.New(db)
	config.db = dbQueries
	log.Println("connected to database")

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("no port provided")
	}

	plexer := http.NewServeMux()

	plexer.HandleFunc("POST /api/data", config.postData)

	server := &http.Server{
		Addr:              port,
		Handler:           plexer,
		ReadHeaderTimeout: 5 * time.Second,
	}
	log.Println("server started, serving at http://localhost:8080")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal("server not started")
	}
}
