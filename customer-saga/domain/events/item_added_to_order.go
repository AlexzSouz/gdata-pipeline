package events

type ItemAddedToOrder struct {
	DomainEvent
	OrderReference string  `json:"order_reference"`
	Name           string  `json:"item_name"`
	Quantity       int     `json:"quantity"`
	Price          float32 `json:"total_price"`
}
