package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/mrckurz/CI-CD-MCM/internal/handler"
	"github.com/mrckurz/CI-CD-MCM/internal/store"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	s := store.NewMemoryStore()
	h := handler.NewHandler(s)

	r := mux.NewRouter()
	h.RegisterRoutes(r)

	fmt.Printf("Product Catalog API listening on :%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
