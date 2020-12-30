package handle

import (
	"collector/viewModel"
	"github.com/kataras/iris/v12"
)

func Collector(ctx iris.Context) {
	var span viewModel.Span
	if err := ctx.ReadJSON(&span); err != nil {
		ctx.StopExecution()
		return
	}
	collect(span)
}

func collect(apm viewModel.Span) {

}
