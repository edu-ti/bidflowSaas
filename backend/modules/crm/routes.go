package crm

import (
	"net/http"
	"github.com/gorilla/mux"
)

// RegisterRoutes creates the routes for the CRM module and attaches them to the provided router.
func RegisterRoutes(router *mux.Router, handler *Handler) {
	if handler == nil {
		return
	}
	router.HandleFunc("/customers", handler.ListCustomers).Methods(http.MethodGet)
	router.HandleFunc("/customers", handler.CreateCustomer).Methods(http.MethodPost)
	router.HandleFunc("/opportunities", handler.ListOpportunities).Methods(http.MethodGet)
}
