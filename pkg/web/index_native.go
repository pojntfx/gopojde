// +build !wasm,!js

package web

import "embed"

//go:embed companion
var CompanionFS embed.FS

//go:embed manager
var ManagerFS embed.FS
