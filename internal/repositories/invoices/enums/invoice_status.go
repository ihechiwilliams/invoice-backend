package enums

// InvoiceStatus ENUM(PENDING_PAYMENT, DRAFT, OVERDUE, PAID)
//
//go:generate go run github.com/abice/go-enum@v0.5.5
type InvoiceStatus string
