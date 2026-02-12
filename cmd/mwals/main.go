package main

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"flag"
	"log"
	"net/http"

	"github.com/frankmwase/malawi-pay-standard/pkg/mwals"
)

func main() {
	port := flag.String("port", "8080", "HTTP port to listen on")
	dataPath := flag.String("data", "als_data.json", "Path to JSON data store")
	keyHex := flag.String("key", "", "Ed25519 private key in hex (optional)")
	flag.Parse()

	var key ed25519.PrivateKey
	var err error

	if *keyHex != "" {
		seed, err := hex.DecodeString(*keyHex)
		if err != nil {
			log.Fatalf("Invalid private key hex: %v", err)
		}
		key = ed25519.NewKeyFromSeed(seed)
	} else {
		log.Println("No signing key provided. Generating a fresh one for this session...")
		_, key, err = ed25519.GenerateKey(rand.Reader)
		if err != nil {
			log.Fatalf("Failed to generate key: %v", err)
		}
		// Print the seed so the user can save it
		seedHex := hex.EncodeToString(key.Seed())
		log.Printf("Session Private Key (Seed): %s", seedHex)
	}

	service, err := mwals.NewService(key, *dataPath)
	if err != nil {
		log.Fatalf("Failed to initialize ALS service: %v", err)
	}

	handler := mwals.NewHandler(service)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", handler.Health)
	mux.HandleFunc("/resolve/", handler.ServeHTTP)
	mux.HandleFunc("/register", handler.Register)

	log.Printf("MW-ALS (The Discovery) starting on port %s...", *port)
	log.Printf("Data store: %s", *dataPath)

	server := &http.Server{
		Addr:    ":" + *port,
		Handler: mux,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
