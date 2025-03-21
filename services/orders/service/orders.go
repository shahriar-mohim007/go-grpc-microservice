package service

import (
	"context"

	"github.com/shahriar-mohim007/kitchen/services/common/genproto/orders"
)

var ordersDb = make([]*orders.Order, 0)

type OrderService struct {
	// store
}

func NewOrderService() *OrderService {
	return &OrderService{}
}

func (s *OrderService) CreateOrder(ctx context.Context, order *orders.Order) error {
	ordersDb = append(ordersDb, order)
	return nil
}

func (s *OrderService) GetOrders(ctx context.Context, customerID int32) []*orders.Order {
	var filteredOrders []*orders.Order

	for _, order := range ordersDb {
		if order.CustomerID == customerID {
			filteredOrders = append(filteredOrders, order)
		}
	}

	return filteredOrders
}
