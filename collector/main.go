package main

func main() {
	collect := New()
	collect.StartServer(&CollectOptions{
		CollectorGRPCHostPort: "127.0.0.1:33666",
	})
	defer collect.Close()
}
