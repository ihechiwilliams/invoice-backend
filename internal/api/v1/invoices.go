package v1

import (
	"context"
	"errors"
	"fmt"
	"invoice-backend/internal/constants"
	"net/http"
	"strconv"
	"strings"

	"invoice-backend/internal/api/server"
	"invoice-backend/internal/repositories/invoices"
	"invoice-backend/internal/repositories/invoices/enums"
	"invoice-backend/internal/repositories/invoicesitems"
	"invoice-backend/internal/shared"

	"github.com/go-chi/render"
	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/samber/lo"
)

type InvoiceHandler struct {
	invoicesRepo      invoices.Repository
	invoicesItemsRepo invoicesitems.Repository
}

func NewInvoiceHandler(invoicesRepo invoices.Repository, invoiceItemsRepo invoicesitems.Repository) *InvoiceHandler {
	return &InvoiceHandler{
		invoicesRepo:      invoicesRepo,
		invoicesItemsRepo: invoiceItemsRepo,
	}
}

func (a *API) V1GetInvoices(w http.ResponseWriter, r *http.Request, reqBody server.V1GetInvoicesParams) {
	var (
		invoiceFilter *invoices.InvoiceDBFilter
		page          = reqBody.Data.Page
		pageSize      = reqBody.Data.PageSize
	)

	fmt.Println("reach here")

	params := reqBody.Data

	if params.Filters != nil {
		filter, prepareErr := prepareInvoiceFilter(lo.FromPtr(params.Filters))
		if prepareErr != nil {
			server.BadRequestError(prepareErr, w, r)
			return
		}

		invoiceFilter = filter
	}

	if page == nil {
		page = getDefaultPage()
	}

	if pageSize == nil {
		pageSize = getDefaultPageSize()
	}

	paginationFilter := preparePagination(pageSize, page)
	fmt.Println(paginationFilter)
	fmt.Println("reach here")

	result, err := a.invoicesHandler.invoicesRepo.ListInvoices(r.Context(), invoiceFilter, paginationFilter)
	if err != nil {
		server.ProcessingError(err, w, r)

		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, lo.Map(result, func(invoice *invoices.Invoice, _ int) server.InvoiceResponseData {
		return serializeInvoiceToAPIResponse(invoice)
	}))
}

func (a *API) V1CreateInvoice(w http.ResponseWriter, r *http.Request) {
	reqBody := new(server.V1CreateInvoiceJSONRequestBody)

	err := render.DecodeJSON(r.Body, reqBody)
	if err != nil {
		server.BadRequestError(err, w, r)

		return
	}

	invoiceData := reqBody.Data
	invoiceNum, err := a.invoicesHandler.GenerateNextInvoiceNumber(r.Context())
	if err != nil {
		server.BadRequestError(err, w, r)
		return
	}
	newInvoice := &invoices.DBInvoice{
		ID:            uuid.New(),
		UserID:        lo.FromPtr(invoiceData.UserId),
		CustomerID:    lo.FromPtr(invoiceData.CustomerId),
		InvoiceNumber: invoiceNum,
		DueDate:       invoiceData.DueDate.Time,
		Status:        enums.InvoiceStatusDRAFT,
	}

	var totalAmount float64

	newInvoice.Items = lo.Map(invoiceData.Items, func(item server.Item, _ int) *invoicesitems.InvoiceItem {
		totalPrice := float64(float32(lo.FromPtr(item.Quantity)) * lo.FromPtr(item.UnitPrice))

		totalAmount += totalPrice

		return &invoicesitems.InvoiceItem{
			ID:          uuid.New(),
			Description: lo.FromPtr(item.Description),
			InvoiceID:   newInvoice.ID,
			Quantity:    lo.FromPtr(item.Quantity),
			UnitPrice:   float64(lo.FromPtr(item.UnitPrice)),
			TotalPrice:  totalPrice,
		}
	})

	newInvoice.TotalAmount = totalAmount

	result, err := a.invoicesHandler.invoicesRepo.CreateInvoice(r.Context(), newInvoice)
	if err != nil {
		server.ProcessingError(err, w, r)

		return
	}

	response := serializeInvoiceToAPIResponse(result)
	render.Status(r, http.StatusCreated)
	render.JSON(w, r, response)
}

func serializeInvoiceToAPIResponse(invoice *invoices.Invoice) server.InvoiceResponseData {
	items := lo.Map(invoice.Items, func(items *invoicesitems.InvoiceItem, _ int) server.Item {
		return serializeInvoiceItemsToAPIResponse(items)
	})

	return server.InvoiceResponseData{
		DueDate:     openapi_types.Date{Time: invoice.DueDate},
		Id:          invoice.ID,
		Items:       items,
		Sender:      "nil",
		Status:      server.InvoiceStatusEnum(invoice.Status),
		TotalAmount: lo.ToPtr(float32(invoice.TotalAmount)),
	}
}

func serializeInvoiceItemsToAPIResponse(item *invoicesitems.InvoiceItem) server.Item {
	return server.Item{
		Description: &item.Description,
		Id:          &item.ID,
		InvoiceId:   &item.InvoiceID,
		Quantity:    &item.Quantity,
		TotalPrice:  lo.ToPtr(float32(item.TotalPrice)),
		UnitPrice:   lo.ToPtr(float32(item.UnitPrice)),
	}
}

func prepareInvoiceFilter(filter server.InvoiceFilters) (*invoices.InvoiceDBFilter, error) {
	validationErr := validateInvoicesFilter(filter)
	if validationErr != nil {
		return nil, validationErr
	}

	ids := make([]*uuid.UUID, 0)
	userIds := make([]*uuid.UUID, 0)
	customerIds := make([]*uuid.UUID, 0)
	invoiceNumbers := make([]*string, 0)
	invoiceStatus := make([]*enums.InvoiceStatus, 0)

	if filter.Id != nil {
		for _, id := range lo.FromPtr(filter.Id) {
			invId := uuid.MustParse(id)
			ids = append(ids, &invId)
		}
	}

	if filter.UserId != nil {
		for _, userID := range lo.FromPtr(filter.UserId) {
			usId := uuid.MustParse(userID)
			userIds = append(userIds, &usId)
		}
	}

	if filter.CustomerId != nil {
		for _, customerID := range lo.FromPtr(filter.CustomerId) {
			cusId := uuid.MustParse(customerID)
			customerIds = append(customerIds, &cusId)
		}
	}

	if filter.InvoiceNumber != nil {
		for _, invoiceNumber := range lo.FromPtr(filter.InvoiceNumber) {
			invoiceNumbers = append(invoiceNumbers, &invoiceNumber)
		}
	}

	if filter.Status != nil {
		for _, stat := range lo.FromPtr(filter.Status) {
			status, parseErr := enums.ParseInvoiceStatus(string(stat))
			if parseErr != nil {
				return nil, fmt.Errorf("invalid invoice status: %s", parseErr.Error())
			}

			invoiceStatus = append(invoiceStatus, lo.ToPtr(status))
		}
	}

	return &invoices.InvoiceDBFilter{
		CustomerID:    customerIds,
		UserID:        userIds,
		ID:            ids,
		InvoiceNumber: invoiceNumbers,
		Status:        invoiceStatus,
	}, nil
}

func validateInvoicesFilter(filter server.InvoiceFilters) error {
	if filter.Id == nil && filter.CustomerId == nil && filter.UserId == nil && filter.InvoiceNumber == nil && filter.Status == nil {
		return errors.New("at least one filter must be provided")
	}

	return nil
}

func preparePagination(limit, page *int) shared.Pagination {
	return shared.Pagination{
		Limit: limit,
		Page:  page,
	}
}

func getDefaultPageSize() *int {
	return lo.ToPtr(constants.DefaultPageSize)
}

func getDefaultPage() *int {
	return lo.ToPtr(constants.DefaultPageNumber)
}

func (h *InvoiceHandler) GenerateNextInvoiceNumber(ctx context.Context) (string, error) {
	lastInvoice, err := h.invoicesRepo.FetchLastInvoice(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to fetch last invoice: %w", err)
	}

	// Default to "INV0000001" if no invoices exist
	if lastInvoice == nil || lastInvoice.InvoiceNumber == "" {
		return "INV0000001", nil
	}

	// Extract the numeric part of the last invoice number
	lastNumber := strings.TrimPrefix(lastInvoice.InvoiceNumber, "INV")
	parsedNumber, parseErr := strconv.Atoi(lastNumber)
	if parseErr != nil {
		return "", fmt.Errorf("invalid invoice number format: %w", parseErr)
	}

	// Increment the numeric part and format it with leading zeros
	nextNumber := fmt.Sprintf("INV%07d", parsedNumber+1)
	return nextNumber, nil
}
