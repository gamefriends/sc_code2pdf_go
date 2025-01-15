//go:build embedfonts
// +build embedfonts

package main

import "embed"

//go:embed fonts/*
var embeddedFonts embed.FS
