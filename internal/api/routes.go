package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	openapi_types "github.com/oapi-codegen/runtime/types"

	"invoice-backend/internal/api/server"
	v1 "invoice-backend/internal/api/v1"
)

func InitRoutes(router *chi.Mux, si *Routes) {
	server.HandlerFromMux(si, router)
}
func NewRoutes(apiV1 *v1.API) *Routes {
	return &Routes{
		v1: apiV1,
	}
}

// Routes is the wrapper for all the versions of the API defined by server.ServerInterface.
type Routes struct {
	v1 *v1.API
}

func (a Routes) V1GetActivities(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (a Routes) V1GetCustomers(w http.ResponseWriter, r *http.Request, params server.V1GetCustomersParams) {
	a.v1.V1GetCustomers(w, r, params)
}

func (a Routes) V1CreateCustomer(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (a Routes) V1GetInvoices(w http.ResponseWriter, r *http.Request, params server.V1GetInvoicesParams) {
	//TODO implement me
	panic("implement me")
}

func (a Routes) V1CreateInvoice(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (a Routes) V1DeleteInvoice(w http.ResponseWriter, r *http.Request, invoiceId string) {
	//TODO implement me
	panic("implement me")
}

func (a Routes) V1GetInvoice(w http.ResponseWriter, r *http.Request, invoiceId openapi_types.UUID) {
	//TODO implement me
	panic("implement me")
}

func (a Routes) V1UpdateInvoice(w http.ResponseWriter, r *http.Request, invoiceId string) {
	//TODO implement me
	panic("implement me")
}
