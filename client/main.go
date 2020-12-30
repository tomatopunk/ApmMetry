package main

import (
	"github.com/kataras/iris/v12"
	"html/template"
)

func main() {
	app := iris.Default()
	app.Handle("GET", "/hello", func(ctx iris.Context) {
		ctx.JSON(iris.Map{"message": "hello!!!"})
	})

	app.Handle("GET", "/api1", func(ctx iris.Context) {
		ctx.JSON(template.OK)
	})

	app.Handle("GET", "/api2", func(ctx iris.Context) {
		ctx.JSON(template.OK)
	})

	app.Handle("GET", "/api3", func(ctx iris.Context) {
		ctx.JSON(template.OK)
	})

	app.Listen(":8080")
}

func api1() {

}

func api2() {

}

func api3() {

}
