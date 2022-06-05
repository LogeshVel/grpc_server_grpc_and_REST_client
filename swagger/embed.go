package OpenAPI

import (
	"embed"
)

//go:embed OpenAPI/*
var OpenAPI embed.FS
