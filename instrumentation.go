package nriris

import (
	"github.com/newrelic/go-agent"

	iris "gopkg.in/kataras/iris.v6"
)

const NewRelicTransaction = "__newrelic_transaction__"

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
