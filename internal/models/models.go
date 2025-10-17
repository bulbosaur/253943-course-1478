package models

type Order struct {
	ID       string
	Item     string
	Quantity int32
}

type CreateOrderRequest struct {
	Item     string
	Quantity int32
}

type CreateOrderResponse struct {
	ID string
}

type DeleteOrderRequest struct {
	ID string
}

type DeleteOrderResponse struct {
	Success bool
}

type GetOrderRequest struct {
	ID string
}

type GetOrderResponse struct {
	Order Order
}

type ListOrdersResponse struct {
	Orders []Order
}

type UpdateOrderRequest struct {
	ID       string
	Item     string
	Quantity int32
}

type UpdateOrderResponse struct {
	Order Order
}