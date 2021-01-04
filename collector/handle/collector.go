package handle

import (
	"collector/produce"
	"github.com/kataras/iris/v12"
)

func Collector(ctx iris.Context) {
	var span produce.Span
	if err := ctx.ReadJSON(&span); err != nil {
		ctx.StopExecution()
		return
	}
	collect(&span)
}

func collect(apm *produce.Span) {
	if err := (produce.Redis{Span: apm}.SendMessage()); err != nil {
		panic(err)
	}
}
