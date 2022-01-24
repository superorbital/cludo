//go:build tools
// +build tools

// Inspired by: https://play-with-go.dev/tools-as-dependencies_go115_en/

package tools

import (
	_ "github.com/ahmetb/govvv"
	_ "github.com/mitchellh/gox"
	_ "github.com/ofabry/go-callvis"
)
