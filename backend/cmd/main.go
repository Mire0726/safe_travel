package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/Mire0726/safe_travel/backend/api/infrastructure"
	"github.com/Mire0726/safe_travel/backend/api/infrastructure/firebase"
	"github.com/Mire0726/safe_travel/backend/cmd/server"
)

func main() {
	// データベース設定の読み込み
	dbConfig, err := infrastructure.LoadDBConfig()
	if err != nil {
		log.Fatal("Could not load database config:", err)
	}

	// データベース接続の初期化
	db, err := infrastructure.NewDB(dbConfig)
	if err != nil {
		log.Fatal("Could not initialize database:", err)
	}

	firebaseAuth, err := firebase.NewFirebaseAuth()
	if err != nil {
		log.Fatalf("Failed to initialize Firebase: %v", err)
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
	server := server.NewServer(db, firebaseAuth)
	server.Serve(":8080")

	if db == nil {
		log.Fatal("Database connection is nil in main")
	}
}
