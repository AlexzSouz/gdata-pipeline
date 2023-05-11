package domain

import (
	"github.com/gdata/customer-saga/domain/aggregates"
	"github.com/gdata/customer-saga/domain/commands"
	"github.com/gdata/customer-saga/domain/entities"
	"github.com/gdata/customer-saga/domain/events"
)

func CreateCustomer(command commands.CreateCustomerCommand) aggregates.ICustomer {
	customer := &aggregates.Customer{
		Aggregate: aggregates.Aggregate{
			Id: command.Id,
		},
		FirstName: command.FirstName,
		LastName:  command.LastName,
		Reference: command.Reference,
		Status:    command.Status,
		Orders:    make(map[string]entities.Order),
	}

	customer.Raise(events.CustomerCreated{
		DomainEvent: events.DomainEvent{
			Id: customer.Id,
		},
		FirstName: command.FirstName,
		LastName:  command.LastName,
		Reference: command.Reference,
		Status:    command.Status,
	})

	return customer
}
