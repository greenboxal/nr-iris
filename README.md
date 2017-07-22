# nr-iris

NewRelic instrumentation helpers for iris

## Usage

###  Instrumenting web transactions

```go
app.HandleFunc("POST", "/orders", nriris.WrapHandler("/orders", createOrder))
```

### Getting newrelic.Transaction:
```go
func handler (ctx *iris.Context) {
    txn := nriris.GetTransaction(ctx)
}
```

For complete documentation, check [here](https://godoc.org/github.com/greenboxal/nr-iris).

## License

See [here](LICENSE).

