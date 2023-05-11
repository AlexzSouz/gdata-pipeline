package enums

type OrderStatus string

const (
	OrderPending   OrderStatus = "Pending"
	OrderDelivered OrderStatus = "Delivered"
	OrderCancelled OrderStatus = "Cancelled"
)
