package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/Mire0726/safe_travel/backend/api/infrastructure"
)

func main() {
	db, err := infrastructure.ConnectToDB()
	if err != nil {
		log.Fatal("Could not initialize database:", err)
	}

	flag.Parse()
	defaultPort := "8080"
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
		flag.StringVar(&port, "addr", defaultPort, "default server port")
	}

	addr := fmt.Sprintf(":%s", port)
	log.Printf("Listening on %s...\n", addr)
	infrastructure.Serve(addr)

	if db == nil {
		log.Fatal("Database connection is nil in main")
	}
}
