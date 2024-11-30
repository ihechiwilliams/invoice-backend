package customers

func FromDBCustomer(dbCustomer *DBCustomer) *Customer {
	return &Customer{
		ID:        dbCustomer.ID,
		Name:      dbCustomer.Name,
		Email:     dbCustomer.Email,
		Phone:     dbCustomer.Phone,
		Address:   dbCustomer.Address,
		UserID:    dbCustomer.UserID,
		CreatedAt: dbCustomer.CreatedAt,
		UpdatedAt: dbCustomer.UpdatedAt,
	}
}
