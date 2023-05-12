package aggregates

import (
	"fmt"
	"strings"

	"github.com/gdata/customer-saga/domain/entities"
	"github.com/gdata/customer-saga/domain/enums"
	"github.com/gdata/customer-saga/domain/events"
)

type ICustomer interface {
	AddOrder(order entities.Order)
	AddLineItemToOrder(orderReference string, lineItem entities.LineItem)
	// CalculateOrderTotal() float32
}

// Customer Aggregate struct following DDD pratices
type Customer struct {
	Aggregate
	FirstName   string
	LastName    string
	Reference   string
	Status      enums.CustomerStatus
	Orders      map[string]entities.Order
	AmountSpent float32
}

func (c *Customer) AddOrder(order entities.Order) {
	if _, found := c.Orders[order.Reference]; found {
		err := fmt.Errorf("Order reference [%q] already present.", order.Reference)
		panic(err)
	}

	if order.Reference == "" || len(strings.TrimSpace(order.Reference)) == 0 {
		err := fmt.Errorf("Invalid order reference [%q].", order.Reference)
		panic(err)
	}

	order.LineItems = make(map[int]entities.LineItem)

	c.Orders[order.Reference] = order

	// Raising domain events. For Event-Sourcing architecture
	c.Raise(events.OrderAddedToCustomer{
		DomainEvent: events.DomainEvent{
			Id: order.Id,
		},
		CustomerReference: c.Reference,
		Reference:         order.Reference,
		Status:            order.Status,
		Timestamp:         order.Timestamp.Unix(),
	})
}

func (c *Customer) AddLineItemToOrder(orderReference string, lineItem entities.LineItem) {
	order := c.Orders[orderReference]
	items := order.LineItems

	if item, found := items[lineItem.Id]; found {
		item.Quantity += lineItem.Quantity
		item.Price += lineItem.Price
		lineItem = item
	}

	items[lineItem.Id] = lineItem
	c.AmountSpent += lineItem.Price

	// Raising domain events. For Event-Sourcing architecture
	c.Raise(events.ItemAddedToOrder{
		DomainEvent: events.DomainEvent{
			Id: order.Id,
		},
		OrderReference: orderReference,
		Name:           lineItem.Name,
		Quantity:       lineItem.Quantity,
		Price:          lineItem.Price,
	})
}
