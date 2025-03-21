package main

func main() {
	grpcServer := NewGRPCServer(":9000")
	grpcServer.Run()
}
