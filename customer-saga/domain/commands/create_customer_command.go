package commands

import . "github.com/gdata/customer-saga/domain/enums"

type CreateCustomerCommand struct {
	Id        int            `json:"id"`
	FirstName string         `json:"first_name"`
	LastName  string         `json:"last_name"`
	Reference string         `json:"customer_reference"`
	Status    CustomerStatus `json:"status"`
}
