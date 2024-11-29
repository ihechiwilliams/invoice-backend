package invoices

import "github.com/samber/lo"

func FromDBInvoice(dbInvoice *DBInvoice) *Invoice {
	return &Invoice{
		ID:            dbInvoice.ID,
		CustomerID:    dbInvoice.CustomerID,
		InvoiceNumber: dbInvoice.InvoiceNumber,
		Status:        dbInvoice.Status,
		TotalAmount:   dbInvoice.TotalAmount,
		DueDate:       dbInvoice.DueDate,
		IssueDate:     dbInvoice.IssueDate,
		Items:         dbInvoice.Items,
		CreatedAt:     dbInvoice.CreatedAt,
		UpdatedAt:     dbInvoice.UpdatedAt,
	}
}

func FromDBInvoiceList(dbInvoices []*DBInvoice) []*Invoice {
	return lo.Map(dbInvoices, func(invoice *DBInvoice, _ int) *Invoice {
		return FromDBInvoice(invoice)
	})
}
