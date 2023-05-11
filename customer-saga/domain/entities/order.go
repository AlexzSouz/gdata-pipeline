package entities

import (
	"time"

	"github.com/gdata/customer-saga/domain/enums"
)

type Order struct {
	Id        int               `json:"id"`
	Reference string            `json:"order_reference"`
	Status    enums.OrderStatus `json:"status"`
	Timestamp time.Time         `json:"order_timestamp"`
	LineItems map[int]LineItem  `json:"line_items"`
	// CustomerReference string            `json:"customer_reference"`
}
