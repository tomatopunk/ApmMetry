package main

import (
	"github.com/kataras/iris/v12"
)

func main() {
	app := iris.Default()
	//app.Handle("POST", "/collect", handle.Collector)
	app.Listen(":8081")
}
