package interop

import (
	"encoding/json"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

func Unmarshal(data interface{}, out interface{}) error {
	return json.Unmarshal([]byte(app.Window().Get("JSON").Call("stringify", data).String()), &out)
}
