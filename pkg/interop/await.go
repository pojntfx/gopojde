package interop

import (
	"errors"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

// See https://stackoverflow.com/questions/68426700/how-to-wait-a-js-async-function-from-golang-wasm
func Await(awaitable app.Value) (app.Value, error) {
	then := make(chan app.Value)
	defer close(then)
	thenFunc := app.FuncOf(func(this app.Value, args []app.Value) interface{} {
		then <- args[0]

		return nil
	})
	defer thenFunc.Release()

	catch := make(chan app.Value)
	defer close(catch)
	catchFunc := app.FuncOf(func(this app.Value, args []app.Value) interface{} {
		catch <- args[0]

		return nil
	})
	defer catchFunc.Release()

	awaitable.Call("then", thenFunc).Call("catch", catchFunc)

	select {
	case result := <-then:
		return result, nil
	case err := <-catch:
		return nil, errors.New(err.String())
	}
}
