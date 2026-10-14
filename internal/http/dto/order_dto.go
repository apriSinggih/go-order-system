package dto

type CreateOrderItemRequest struct {
	ProductID int64 `json:"product_id" binding:"required"`
	Quantity  int   `json:"quantity" binding:"required,gt=0"`
}

type CreateOrderRequest struct {
	UserID int64                    `json:"user_id" binding:"required"`
	Items  []CreateOrderItemRequest `json:"items" binding:"required,dive"`
}