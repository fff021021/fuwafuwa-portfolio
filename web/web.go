package web

import "embed"

// EmbeddedFiles は静的ファイルとindex.htmlをバイナリ内に埋め込みます
//
//go:embed all:static index.html
var EmbeddedFiles embed.FS
