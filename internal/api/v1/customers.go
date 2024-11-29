package v1

import (
	"invoice-backend/internal/repositories/customers"
	"net/http"
)

type CustomersHandler struct {
	customersRepo customers.Repository
}

func NewCustomersHandler(customersRepo customers.Repository) *CustomersHandler {
	return &CustomersHandler{
		customersRepo: customersRepo,
	}
}

func (a *API) V1GetCustomers(w http.ResponseWriter, r *http.Request)
