package domain

import (
	"time"

	"github.com/asaskevich/govalidator"
)

type OrderItem struct {
	ProductID string  `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

type Order struct {
	ID                   string      `json:"id"`
	UserID               string      `json:"user_id" gorm:"type:uuid;notnull"`
	OrderDate            time.Time   `json:"order_date"`
	Status               string      `json:"status"`
	OrderItems           []OrderItem `json:"order_items"`
	DeliveryAddress      string      `json:"delivery_address"`
	PaymentMethod        string      `json:"payment_method"`
	TotalCost            float64     `json:"total_cost"`
	CreatedAt            time.Time   `json:"created_at"`
	ExpectedDeliveryDate time.Time   `json:"expected_delivery_date"`
	SellerID             string      `json:"seller_id"`
}

func init() {
	govalidator.SetFieldsRequiredByDefault(true)
}

func NewOrder() *Order {
	return &Order{}

}

func (order *Order) Validate() error {

	_, err := govalidator.ValidateStruct(order)

	if err != nil {
		return err
	}

	return nil
}

func (o *Order) AddItem(item OrderItem) {
	o.OrderItems = append(o.OrderItems, item)
}

func (o *Order) prepare() {
	o.TotalCost = 0
	for _, item := range o.OrderItems {
		o.TotalCost += item.Price * float64(item.Quantity)
	}
	o.OrderDate = time.Now()
	o.CreatedAt = time.Now()
}
