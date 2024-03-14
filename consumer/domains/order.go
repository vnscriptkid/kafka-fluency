package domains

type OrderCreated struct {
	OrderID string `json:"orderID"`
	Item    string `json:"item"`
}
