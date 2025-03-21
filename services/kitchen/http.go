package main

import (
	"context"
	"log"
	"net/http"
	"text/template"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/shahriar-mohim007/kitchen/services/common/genproto/orders"
)

type HttpServer struct {
	addr string
}

func NewHttpServer(addr string) *HttpServer {
	return &HttpServer{addr: addr}
}

func (s *HttpServer) Run() error {
	router := http.NewServeMux()

	conn := NewGRPCClient(":9000")
	defer conn.Close()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		c := orders.NewOrderServiceClient(conn)

		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		_, err := c.CreateOrder(ctx, &orders.CreateOrderRequest{
			CustomerID: 24,
			ProductID:  3123,
			Quantity:   2,
		})
		if err != nil {
			log.Printf("gRPC error: %v", err)
			http.Error(w, "Failed to create order", http.StatusInternalServerError)
			return
		}

		res, err := c.GetOrders(ctx, &orders.GetOrdersRequest{
			CustomerID: 42,
		})
		if err != nil {
			log.Printf("gRPC error: %v", err)
			http.Error(w, "Failed to fetch orders", http.StatusInternalServerError)
			return
		}

		t := template.Must(template.New("orders").Parse(ordersTemplate))
		if err := t.Execute(w, res.GetOrders()); err != nil {
			log.Printf("template error: %v", err)
			http.Error(w, "Failed to render template", http.StatusInternalServerError)
			return
		}
	})

	log.Println("Starting HTTP server on", s.addr)
	return http.ListenAndServe(s.addr, router)
}

func NewGRPCClient(addr string) *grpc.ClientConn {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("gRPC connection failed: %v", err)
	}
	return conn
}

var ordersTemplate = `
<!DOCTYPE html>
<html>
<head>
    <title>Kitchen Orders</title>
</head>
<body>
    <h1>Orders List</h1>
    <table border="1">
        <tr>
            <th>Order ID</th>
            <th>Customer ID</th>
            <th>Quantity</th>
        </tr>
        {{range .}}
        <tr>
            <td>{{.OrderID}}</td>
            <td>{{.CustomerID}}</td>
            <td>{{.Quantity}}</td>
        </tr>
        {{end}}
    </table>
</body>
</html>`
