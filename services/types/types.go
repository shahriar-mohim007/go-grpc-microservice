package types

import (
	"context"

	"github.com/shahriar-mohim007/kitchen/services/common/genproto/orders"
)

type OrderService interface {
	CreateOrder(context.Context, *orders.Order) error
	GetOrders(ctx context.Context, customerID int32) []*orders.Order
}
