package licitacoes

import (
	"net/http"
	"github.com/gorilla/mux"
)

// RegisterRoutes creates the routes for the licitacoes module and attaches them to the provided router.
func RegisterRoutes(router *mux.Router, handler *Handler) {
	if handler == nil {
		return
	}
	router.HandleFunc("/editais", handler.ListEditais).Methods(http.MethodGet)
	router.HandleFunc("/editais", handler.CreateEdital).Methods(http.MethodPost)
}
