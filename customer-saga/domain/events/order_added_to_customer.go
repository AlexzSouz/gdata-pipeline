package events

import (
	"github.com/gdata/customer-saga/domain/enums"
)

type OrderAddedToCustomer struct {
	DomainEvent
	CustomerReference string            `json:"customer_reference"`
	Reference         string            `json:"order_reference"`
	Status            enums.OrderStatus `json:"status"`
	Timestamp         int64             `json:"order_timestamp"`
}
