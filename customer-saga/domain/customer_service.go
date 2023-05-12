package domain

import (
	"encoding/csv"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gdata/customer-saga/abstractions"
	"github.com/gdata/customer-saga/domain/aggregates"
	"github.com/gdata/customer-saga/domain/commands"
	"github.com/gdata/customer-saga/domain/entities"
	"github.com/gdata/customer-saga/domain/enums"
	"github.com/gdata/customer-saga/messages"
)

type CustomerService struct {
}

type CsvEntry struct {
	Type    string
	Headers map[string]int
	Entries [][]string
}

func (c *CustomerService) Process(ctx abstractions.IAppContext, path string) {
	ctx, span := ctx.CreateSpan("Process")
	defer span.End()

	if strings.HasSuffix(path, ".lock.csv") {
		return
	}

	dirPath, entityToken, dateToken, err := retrievePathTokens(path)
	if err != nil {
		ctx.Logger().Println(err)
		return
	}
	if entityToken != "customers" {
		ctx.Logger().Printf("Waiting for '/%v_%v.csv' file to start processing.", entityToken, dateToken)
		return
	}
	defer func() {
		if err := recover(); err != nil {
			// TODO : Eitheer unlock file, retry, or push error message...
		}
	}()

	newFileName, err := lockCsvFile(dirPath, fmt.Sprintf("%v_%v", entityToken, dateToken))
	if err != nil {
		ctx.Logger().Printf("Failed to lock file. %q", err)
		return
	}

	files := map[string]chan [][]string{
		"customers": make(chan [][]string),
		"orders":    make(chan [][]string),
		"items":     make(chan [][]string),
	}

	for key, cn := range files {
		file := fmt.Sprintf("%v_%v.csv", key, dateToken)
		if key == "customers" {
			file = newFileName
		}
		path := fmt.Sprintf("%v/%v", dirPath, file)
		go retrieveCsvEntries(ctx, path, cn)
	}

	csvEntries := make(map[string]CsvEntry, len(files))
	for k, c := range files {
		data := <-c
		csvEntries[k] = CsvEntry{
			Type:    k,
			Headers: createHeadersMap(data[0]),
			Entries: data[1:],
		}
	}

	customerAggregates := processCustomers(ctx, csvEntries)

	for _, a := range customerAggregates {
		aggregate := a.(*aggregates.Customer)

		event := messages.Envelope[any]{
			TraceId: ctx.GetTraceId(),
			SpanId:  ctx.GetSpanId(),
			Type:    fmt.Sprintf("%T", a),
			Payload: struct {
				Type        string  `json:"type"`
				Reference   string  `json:"customer_reference"`
				OrdersCount int     `json:"number_of_orders"`
				Total       float32 `json:"total_amount_spent"`
			}{
				Type:        fmt.Sprintf("%T", aggregate),
				Reference:   aggregate.Reference,
				OrdersCount: len(aggregate.Orders),
				Total:       aggregate.AmountSpent,
			},
		}

		// TODO : Send external events
		//		 -> bus.Send(event)
		ctx.Logger().Println(event)
	}
}

func processCustomers(ctx abstractions.IAppContext, csvEntries map[string]CsvEntry) []aggregates.ICustomer {
	customers := csvEntries["customers"]
	customersChannels := []chan aggregates.ICustomer{}

	for i, entry := range customers.Entries {
		doneChannel := make(chan aggregates.ICustomer)
		customersChannels = append(customersChannels, doneChannel)

		processCustomer := func(index int, payload []string, done chan aggregates.ICustomer) {
			ctx.Logger().Printf("Processign entry [%v].", index)
			customerId, err := strconv.Atoi(entry[customers.Headers["id"]])
			if err != nil {
				panic("TODO : Treat error here")
			}

			command := commands.CreateCustomerCommand{
				Id:        customerId,
				Reference: strings.TrimSpace(entry[customers.Headers["customer_reference"]]),
				FirstName: strings.TrimSpace(entry[customers.Headers["first_name"]]),
				LastName:  strings.TrimSpace(entry[customers.Headers["last_name"]]),
				Status:    enums.CustomerStatus(strings.TrimSpace(entry[customers.Headers["status"]])), // TODO : Add conversion validation
			}

			ctx.Logger().Print("Creating customer aggregate.")
			customerAggregate := CreateCustomer(command)

			processOrders(ctx, customerAggregate, csvEntries["orders"])
			processOrderItems(ctx, customerAggregate, csvEntries["items"])

			ctx.Logger().Printf("Successfully created customer aggregate with ID [%v].", customerAggregate.(*aggregates.Customer).Id)
			// done <- true
			done <- customerAggregate
		}

		go processCustomer(i, entry, doneChannel)
	}

	customersAggregates := make([]aggregates.ICustomer, 0)
	for _, c := range customersChannels {
		customerAggregate := <-c
		ctx.Logger().Printf("Completed processing customer with ID [%v].", customerAggregate.(*aggregates.Customer).Id)
		customersAggregates = append(customersAggregates, customerAggregate)
	}

	return customersAggregates
}

func processOrders(ctx abstractions.IAppContext, customerAggregate aggregates.ICustomer, ordersCsvEntries CsvEntry) {
	orderHeaders := ordersCsvEntries.Headers
	orderEntries := ordersCsvEntries.Entries

	for _, order := range orderEntries {
		customer, ok := customerAggregate.(*aggregates.Customer)
		if !ok {
			panic("TODO : Treat error here")
		}

		orderCustomerRef := strings.TrimSpace(order[orderHeaders["customer_reference"]])
		if orderCustomerRef == customer.Reference {
			orderId, err := strconv.Atoi(order[orderHeaders["id"]])
			if err != nil {
				panic("TODO : Treat error here")
			}
			timestamp, err := strconv.ParseInt(order[orderHeaders["order_timestamp"]], 10, 64)
			if err != nil {
				panic("TODO : Treat error here")
			}

			customerAggregate.AddOrder(entities.Order{
				Id:        orderId,
				Reference: strings.TrimSpace(order[orderHeaders["order_reference"]]),
				Status:    enums.OrderStatus(strings.TrimSpace(order[orderHeaders["order_status"]])),
				Timestamp: time.Unix(timestamp, 0),
			})
			ctx.Logger().Printf("Processing order %q.", strings.TrimSpace(order[orderHeaders["order_reference"]]))
		}
	}
}

func processOrderItems(ctx abstractions.IAppContext, customerAggregate aggregates.ICustomer, lineItemsCsvEntries CsvEntry) {
	lineItemHeaders := lineItemsCsvEntries.Headers
	lineItemEntries := lineItemsCsvEntries.Entries

	for _, lineItem := range lineItemEntries {
		customer, ok := customerAggregate.(*aggregates.Customer)
		if !ok {
			panic("TODO : Treat error here")
		}

		lineItemOrderRef := strings.TrimSpace(lineItem[lineItemHeaders["order_reference"]])
		if _, found := customer.Orders[lineItemOrderRef]; !found {
			return
		}

		lineItemId, err := strconv.Atoi(lineItem[lineItemHeaders["id"]])
		if err != nil {
			panic("TODO : Treat error here")
		}

		quantity, err := strconv.Atoi(lineItem[lineItemHeaders["quantity"]])
		if err != nil {
			panic("TODO : Treat error here")
		}

		price, err := strconv.ParseFloat(lineItem[lineItemHeaders["total_price"]], 32)
		if err != nil {
			panic("TODO : Treat error here")
		}

		customerAggregate.AddLineItemToOrder(lineItemOrderRef, entities.LineItem{
			Id:       lineItemId,
			Name:     strings.TrimSpace(lineItem[lineItemHeaders["item_name"]]),
			Quantity: quantity,
			Price:    float32(price),
		})
	}
}

func retrieveCsvEntries(ctx abstractions.IAppContext, path string, c chan [][]string) {
	ctx.Logger().Printf("Starting the retrieval of Csv entries from %q.", path)

	file, err := os.Open(path)
	if err != nil {
		ctx.Logger().Fatalf("Failed to open file %q.", path)
		ctx.Terminate()
	}

	entries, err := csv.NewReader(file).ReadAll()
	if err != nil {
		ctx.Logger().Fatalf("Failed to read CSV file %q.", path)
		ctx.Terminate()
	}

	fmt.Println(entries)
	ctx.Logger().Printf("Successfully retrieved the Csv entries from %q.", path)
	c <- entries
}

func createHeadersMap(headers []string) map[string]int {
	hm := make(map[string]int)

	for i, h := range headers {
		hm[strings.TrimSpace(h)] = i
	}

	return hm
}

func retrievePathTokens(path string) (dirPath string, entity string, data string, err error) {
	regex := regexp.MustCompile(`^(?P<path>[\w\W]*)/(?P<entity>[\w\W]*)_(?P<date>\d{6,8}).csv$`)
	names := regex.SubexpNames()
	matches := regex.FindStringSubmatch(strings.TrimSpace(path))
	groups := make(map[string]string)

	for i, n := range names {
		groups[n] = matches[i]
	}

	if len(strings.TrimSpace(groups["date"])) == 0 {
		err = fmt.Errorf("Couldn't extract date token from file path %q.", path)
	}

	dirPath = groups["path"]
	entity = groups["entity"]
	data = groups["date"]

	return
}

func lockCsvFile(dirPath string, fileName string) (newFileName string, err error) {
	oldFile := fmt.Sprintf("%v/%v.csv", dirPath, fileName)

	newFileName = fmt.Sprintf("%v.lock.csv", fileName)
	newFile := fmt.Sprintf("%v/%v", dirPath, newFileName)

	err = os.Rename(oldFile, newFile)
	return
}
