package v1

type API struct {
	activitiesHandler *ActivitiesHandler
	customersHandler  *CustomersHandler
	invoicesHandler   *InvoiceHandler
}

func NewAPI(
	activitiesHandler *ActivitiesHandler,
	customersHandler *CustomersHandler,
	invoicesHandler *InvoiceHandler,
) *API {
	return &API{
		activitiesHandler: activitiesHandler,
		customersHandler:  customersHandler,
		invoicesHandler:   invoicesHandler,
	}
}
