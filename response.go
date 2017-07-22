package nriris

import (
	"net/http"

	nrhttp "github.com/newrelic/go-agent/http"

	iris "gopkg.in/kataras/iris.v6"
)

type Response struct {
	res iris.ResponseWriter
	req *http.Request
}

func NewResponse(ctx *iris.Context) Response {
	return Response{
		res: ctx.ResponseWriter,
		req: ctx.Request,
	}
}

func (r Response) Header() nrhttp.Header {
	return r.res.Header()
}

func (r Response) Code() int {
	return r.res.StatusCode()
}

func (r Response) Request() nrhttp.Request {
	return nrhttp.RequestWrapper{r.req}
}
