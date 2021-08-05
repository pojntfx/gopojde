// +build !wasm,!js

package manager

import "embed"

//go:embed assets
var FS embed.FS
