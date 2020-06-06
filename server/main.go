package main

import (
	"github.com/KouT127/billing-slack-notifier/config"
	"github.com/KouT127/billing-slack-notifier/handler"
	"github.com/KouT127/billing-slack-notifier/module"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"log"
	"net/http"
	"os"
)

func main() {
	config.Configure()

	m, err := module.NewSecretManager()
	if err != nil {
		log.Fatalf("%v", err)
	}

	token, err := m.AccessSecret("SLACK_TOKEN")
	if err != nil {
		log.Fatalf("%v", err)
	}
	channelID, err := m.AccessSecret("CHANNEL_ID")
	if err != nil {
		log.Fatalf("%v", err)
	}

	bigQueryClient, err := module.NewBigQueryClient()
	if err != nil {
		log.Fatalf("%v", err)
	}
	slackClient := module.NewSlackClient(token, channelID)
	h := handler.NewHandler(slackClient, bigQueryClient)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Post("/notification", h.NotificationHandler)
	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
