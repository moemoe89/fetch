//go:build tools
// +build tools

// This package imports things required by build/lint scripts, to force `go mod` to see them as dependencies.
// For more details, see https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module
package tools

import (
	_ "github.com/golang/mock/mockgen"
)
