package main

import (
	"context"
	"encoding/json"
	"github.com/shahriar-mohim007/kitchen/services/common/genproto/orders"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"
	"strconv"
	"text/template"
	"time"
)

type CreateOrderRequest struct {
	CustomerID int32 `json:"customer_id"`
	ProductID  int32 `json:"product_id"`
	Quantity   int32 `json:"quantity"`
}

type GetOrdersRequest struct {
	CustomerID int32 `json:"customer_id"`
}
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
	c := orders.NewOrderServiceClient(conn)

	router.HandleFunc("/create-order", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		var reqBody CreateOrderRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		_, err := c.CreateOrder(ctx, &orders.CreateOrderRequest{
			CustomerID: reqBody.CustomerID,
			ProductID:  reqBody.ProductID,
			Quantity:   reqBody.Quantity,
		})

		if err != nil {
			log.Printf("gRPC error: %v", err)
			http.Error(w, "Failed to create order", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"status": "success", "message": "Order created"}`))
	})

	router.HandleFunc("/get-orders", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		customerIDStr := r.URL.Query().Get("customer_id")
		if customerIDStr == "" {
			http.Error(w, "Missing customer_id parameter", http.StatusBadRequest)
			return
		}

		customerID, err := strconv.Atoi(customerIDStr)
		if err != nil {
			http.Error(w, "Invalid customer_id", http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		res, err := c.GetOrders(ctx, &orders.GetOrdersRequest{
			CustomerID: int32(customerID),
		})
		if err != nil {
			log.Printf("gRPC error: %v", err)
			http.Error(w, "Failed to fetch orders", http.StatusInternalServerError)
			return
		}

		t := template.Must(template.New("orders").Parse(ordersTemplate))
		if err := t.Execute(w, res.GetOrders()); err != nil {
			log.Printf("Template error: %v", err)
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
