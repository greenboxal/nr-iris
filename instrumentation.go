package nriris

import (
	"github.com/newrelic/go-agent"
	"net/http"

	iris "gopkg.in/kataras/iris.v6"
)

const NewRelicTransaction = "__newrelic_transaction__"

type NewRelic struct {
	App newrelic.Application
}

func Apply(nr newrelic.Application, app iris.Framework) {
	app.Adapt(iris.Policies{
		RouterWrapperPolicy: func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
			txn := nr.StartTransaction("*", w, r)

			defer func() {
				err := recover()

				if err != nil {
					switch err := err.(type) {
					case error:
						txn.NoticeError(err)
					default:
						txn.NoticeError(errWrapper{err})
					}
				}

				txn.End()
			}()

			next(txn, r)
		},
	})

	app.UseFunc(func(ctx *iris.Context) {
		txn := ctx.ResponseWriter.(newrelic.Transaction)

		txn.SetName(ctx.Request.URL.Path)

		ctx.Set(NewRelicTransaction, txn)

		ctx.Next()
	})
}

func WrapHandlerFunc(app newrelic.Application, name string, handler iris.HandlerFunc) iris.HandlerFunc {
	if app == nil {
		return handler
	}

	return func(ctx *iris.Context) {
		txn := app.StartTransaction(name, nil, ctx.Request)
		defer txn.End()

		defer func() {
			err := recover()

			if err != nil {
				switch err := err.(type) {
				case error:
					txn.NoticeError(err)
				default:
					txn.NoticeError(errWrapper{err})
				}
			} else {
				txn.ResponseSent(NewResponse(ctx))
			}

			txn.End()

			if err != nil {
				panic(err)
			}
		}()

		ctx.Set(NewRelicTransaction, txn)

		handler(ctx)
	}
}

func GetTransaction(ctx *iris.Context) newrelic.Transaction {
	val := ctx.Get(NewRelicTransaction)

	if val == nil {
		return nil
	}

	return val.(newrelic.Transaction)
}
