// Package web holds the site's templates and static assets. They're embedded
// in the binary for `build`; `serve -dev` reads them from disk instead, so
// template and CSS edits show up without recompiling.
package web

import "embed"

//go:embed templates static
var FS embed.FS
