package main

import "github.com/kataras/iris/v12"

func main() {
	app := iris.Default()
	app.Handle("GET", "/hello", func(ctx iris.Context) {
		ctx.JSON(iris.Map{"message": "hello!!!"})
	})
	app.Listen(":8080")
}
