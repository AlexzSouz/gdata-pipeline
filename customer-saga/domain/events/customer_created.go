package events

import "github.com/gdata/customer-saga/domain/enums"

type CustomerCreated struct {
	DomainEvent
	FirstName string               `json:"first_name"`
	LastName  string               `json:"last_name"`
	Reference string               `json:"customer_reference"`
	Status    enums.CustomerStatus `json:"status"`
}
