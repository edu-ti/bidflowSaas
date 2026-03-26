package ai

import (
	"net/http"
	"github.com/gorilla/mux"
)

// RegisterRoutes creates the routes for the AI module and attaches them to the provided router.
func RegisterRoutes(router *mux.Router, handler *Handler) {
	if handler == nil {
		return
	}
	router.HandleFunc("/process-edital", handler.ProcessEdital).Methods(http.MethodPost)
}
