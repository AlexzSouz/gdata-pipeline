package domain

import (
	"reflect"
	"testing"
	"time"

	"github.com/gdata/customer-saga/domain/aggregates"
	"github.com/gdata/customer-saga/domain/commands"
	"github.com/gdata/customer-saga/domain/entities"
	"github.com/gdata/customer-saga/domain/enums"
	"github.com/sciensoft/fluenttests/fluent/contracts"
	"github.com/sciensoft/fluenttests/fluent/integers"
)

func TestFactoryShouldCreateCustomer(t *testing.T) {
	// Arrange
	fluent := contracts.Fluent[aggregates.ICustomer](t)
	expectedType := reflect.TypeOf(&aggregates.Customer{})
	event := commands.CreateCustomerCommand{
		Id:        1,
		FirstName: "Alex",
		LastName:  "Richard",
		Reference: "ade3-11ed",
		Status:    enums.CustomerActive,
	}

	// Act
	customer := CreateCustomer(event)

	// Assert
	fluent.It(customer).
		Should().BeOfType(expectedType).
		And().HaveFieldWithTag("Id", "json").
		And().HaveFieldWithTag("FirstName", "json").
		And().HaveFieldWithTag("LastName", "json").
		And().HaveFieldWithTag("Reference", "json").
		And().HaveFieldWithTag("Status", "json").
		And().HaveFieldWithTag("Orders", "json")
}

func TestFactoryShouldCreateCustomerWithOrders(t *testing.T) {
	// Arrange
	fluent := contracts.Fluent[map[string]entities.Order](t)
	expectedType := reflect.TypeOf(map[string]entities.Order{})
	customer := CreateCustomer(commands.CreateCustomerCommand{
		Id:        1,
		FirstName: "Alex",
		LastName:  "Richard",
		Reference: "ade3-11ed",
		Status:    enums.CustomerActive,
	})
	orders := customer.(*aggregates.Customer).Orders

	// Act
	customer.AddOrder(entities.Order{
		Id:        1,
		Reference: "dc0aa69c",
		Status:    enums.OrderDelivered,
		Timestamp: time.Now(),
	})

	customer.AddOrder(entities.Order{
		Id:        2,
		Reference: "dc0ab1aa",
		Status:    enums.OrderDelivered,
		Timestamp: time.Now(),
	})

	// Assert
	fluent.It(orders).
		Should().BeOfType(expectedType)
}

func TestFactoryShouldCreateCustomerWithOrderAndLineItems(t *testing.T) {
	// Arrange
	fluent := integers.Fluent[int](t)
	customer := CreateCustomer(commands.CreateCustomerCommand{
		Id:        1,
		FirstName: "Alex",
		LastName:  "Richard",
		Reference: "ade3-11ed",
		Status:    enums.CustomerActive,
	})
	order := entities.Order{
		Id:        1,
		Reference: "dc0aa69c",
		Status:    enums.OrderDelivered,
		Timestamp: time.Now(),
	}
	orders := customer.(*aggregates.Customer).Orders

	// Act
	customer.AddOrder(order)
	customer.AddLineItemToOrder(order.Reference, entities.LineItem{
		Id:       1,
		Name:     "Item 1",
		Quantity: 2,
		Price:    10.59,
	})
	customer.AddLineItemToOrder(order.Reference, entities.LineItem{
		Id:       1,
		Name:     "Item 1",
		Quantity: 3,
		Price:    10.59,
	})
	customer.AddLineItemToOrder(order.Reference, entities.LineItem{
		Id:       2,
		Name:     "Item 2",
		Quantity: 5,
		Price:    45.99,
	})

	// Assert
	fluent.It(len(orders)).
		Should().BeGreaterThanOrEqualTo(1)
}
