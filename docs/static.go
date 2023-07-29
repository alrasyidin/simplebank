package docs

import "embed"

//go:embed swagger/*
var StaticSwagger embed.FS
