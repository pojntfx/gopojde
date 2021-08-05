// +build !wasm,!js

package companion

import "embed"

//go:embed assets
var FS embed.FS
