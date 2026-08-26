package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"example.com/familyitinerary/internal/advisor"
	"example.com/familyitinerary/internal/config"
	"example.com/familyitinerary/internal/httpapi"
	"example.com/familyitinerary/internal/store"
)

func main() {
	defaults := config.Default()
	databasePath := flag.String("db", defaults.DatabasePath, "SQLite database path")
	address := flag.String("address", defaults.Address, "HTTP listen address")
	flag.Parse()
	settings := config.FromValues(*databasePath, *address, false)
	if !settings.Valid() {
		log.Fatal("invalid configuration")
	}
	repository, err := store.Open(settings.DatabasePath)
	if err != nil {
		log.Fatal(err)
	}
	defer repository.Close()
	intake := advisor.NewIntakeService(repository)
	conversation := advisor.NewConversationService(repository)
	server := httpapi.New(intake, conversation)
	fmt.Printf("family itinerary advisor listening on %s\n", settings.Address)
	if err := http.ListenAndServe(settings.Address, server.Handler()); err != nil {
		log.Fatal(err)
	}
}
