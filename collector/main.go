package main

func main() {
	collect := CreateCollector()
	err := collect.StartServer(&CollectOptions{
		CollectorGRPCHostPort: ":5000",
	})
	if err != nil {
		println(err.Error())
	}
	defer collect.Close()

}
