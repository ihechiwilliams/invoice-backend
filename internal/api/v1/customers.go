package v1

import (
	"context"
	"github.com/google/uuid"
	"invoice-backend/internal/api/server"
	"invoice-backend/internal/repositories/customers"
	"net/http"

	"github.com/go-chi/render"
	"github.com/samber/lo"
)

type CustomersHandler struct {
	customersRepo customers.Repository
}

func NewCustomersHandler(customersRepo customers.Repository) *CustomersHandler {
	return &CustomersHandler{
		customersRepo: customersRepo,
	}
}

func (h *CustomersHandler) GetCustomers(ctx context.Context, filter *customers.CustomerDBFilter) ([]server.CustomerResponseData, error) {
	list, fetchCustomerErr := h.customersRepo.ListCustomers(ctx, filter)
	if fetchCustomerErr != nil {
		return nil, fetchCustomerErr
	}

	return lo.Map(list, func(customer *customers.Customer, _ int) server.CustomerResponseData {
		return serializeCustomerToAPIResponse(customer)
	}), nil
}

func (a *API) V1GetCustomers(w http.ResponseWriter, r *http.Request, params server.V1GetCustomersParams) {
	var customerFilter *customers.CustomerDBFilter

	if params.Data != nil && params.Data.Filters != nil {
		filter := prepareCustomerFilter(*params.Data.Filters)
		customerFilter = filter
	}

	customerList, fetchCustomerErr := a.customersHandler.GetCustomers(r.Context(), customerFilter)
	if fetchCustomerErr != nil {
		server.ProcessingError(fetchCustomerErr, w, r)

		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, server.CustomersResponse{Data: customerList})
}

func (a *API) V1CreateCustomer(w http.ResponseWriter, r *http.Request) {
	reqBody := new(server.V1CreateCustomerJSONRequestBody)

	err := render.DecodeJSON(r.Body, reqBody)
	if err != nil {
		server.BadRequestError(err, w, r)
		return
	}

	customerData := reqBody.Data

	newCustomer := &customers.DBCustomer{
		ID:      uuid.New(),
		UserID:  customerData.UserId,
		Name:    customerData.Name,
		Email:   customerData.Email,
		Phone:   customerData.Phone,
		Address: customerData.Address,
	}

	result, err := a.customersHandler.customersRepo.CreateCustomer(r.Context(), newCustomer)
	if err != nil {
		server.BadRequestError(err, w, r)
		return
	}

	response := serializeCustomerToAPIResponse(result)
	render.Status(r, http.StatusCreated)
	render.JSON(w, r, response)
}

func prepareCustomerFilter(filters server.CustomerFilters) *customers.CustomerDBFilter {
	return &customers.CustomerDBFilter{
		UserID: lo.FromPtrOr(filters.UserId, []string{}),
	}
}

func serializeCustomerToAPIResponse(customer *customers.Customer) server.CustomerResponseData {
	return server.CustomerResponseData{
		Email: customer.Email,
		Id:    customer.ID,
		Name:  customer.Name,
		Phone: customer.Phone,
	}
}
