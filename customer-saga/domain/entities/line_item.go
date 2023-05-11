package entities

type LineItem struct {
	Id       int     `json:"id"`
	Name     string  `json:"item_name"`
	Quantity int     `json:"quantity"`
	Price    float32 `json:"total_price"`
	// OrderReference string  `json:"order_reference"`
}
