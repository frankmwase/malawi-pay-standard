package main

import (
	"crypto/ed25519"
	"encoding/hex"
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/frankmwase/malawi-pay-standard/pkg/mwals"
)

func main() {
	listen := flag.String("listen", "127.0.0.1:8080", "HTTP bind address (loopback by default)")
	dataPath := flag.String("data", "als_data.json", "Path to JSON data store")
	flag.Parse()

	// Provision secrets out of band. Never print seeds or pass them on the command line.
	seed, err := hex.DecodeString(os.Getenv("MW_ALS_SIGNING_SEED"))
	if err != nil || len(seed) != ed25519.SeedSize {
		log.Fatal("MW_ALS_SIGNING_SEED must contain a 32-byte Ed25519 seed in hex")
	}
	key := ed25519.NewKeyFromSeed(seed)

	service, err := mwals.NewService(key, *dataPath)
	if err != nil {
		log.Fatalf("Failed to initialize ALS service: %v", err)
	}

	handler := mwals.NewAuthenticatedHandler(service, os.Getenv("MW_ALS_REGISTRATION_TOKEN"))

	mux := http.NewServeMux()
	mux.HandleFunc("/health", handler.Health)
	mux.HandleFunc("/resolve/", handler.ServeHTTP)
	// Registration is disabled unless an operator provisions a bearer token.
	if os.Getenv("MW_ALS_REGISTRATION_TOKEN") != "" {
		mux.HandleFunc("/register", handler.Register)
	}

	log.Printf("MW-ALS listening on %s", *listen)
	server := &http.Server{
		Addr:              *listen,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       30 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
