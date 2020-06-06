package handler

import (
	"fmt"
	"github.com/KouT127/billing-slack-notifier/module"
	"log"
	"net/http"
)

type Handler struct {
	*module.SlackClient
	*module.BigQueryClient
}

func NewHandler(sc *module.SlackClient, bc *module.BigQueryClient) *Handler {
	return &Handler{
		SlackClient:    sc,
		BigQueryClient: bc,
	}
}

func (h *Handler) NotificationHandler(w http.ResponseWriter, r *http.Request) {
	results := h.BigQueryClient.FindBill()
	for _, result := range results {
		err := h.SlackClient.NotifyMessage(result)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w)
			return
		}
	}

	log.Println("successful notification to slack")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w)
}
