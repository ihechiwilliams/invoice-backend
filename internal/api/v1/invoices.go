package v1

import "invoice-backend/internal/repositories/invoices"

type InvoiceHandler struct {
	invoiceRepo invoices.Repository
}

func NewInvoiceHandler(invoiceRepo invoices.Repository) *InvoiceHandler {
	return &InvoiceHandler{
		invoiceRepo: invoiceRepo,
	}
}
