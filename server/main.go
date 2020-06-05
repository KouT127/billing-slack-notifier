package main

import (
	"fmt"
	"github.com/KouT127/billing-slack-notifier/config"
	"github.com/KouT127/billing-slack-notifier/module"
	"log"
	"net/http"
	"os"
)

func main() {
	config.Configure()
	http.HandleFunc("/notification", notificationHandler)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}
	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func notificationHandler(w http.ResponseWriter, r *http.Request) {
	slackClient := module.NewSlackClient()
	bigQueryClient := module.NewBigQueryClient()

	results := bigQueryClient.FindBill()
	for _, result := range results {
		err := slackClient.NotifyMessage(result)
		if err != nil {
			log.Fatalf("%v", err)
		}
	}

	log.Println("successful notification to slack")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "ok")
	return
}
