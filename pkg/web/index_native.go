// +build !wasm,!js

package web

import "embed"

//go:embed manager
var ManagerFS embed.FS
