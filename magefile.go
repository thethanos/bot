//go:build mage

package main

import (
	"github.com/magefile/mage/sh"
)

func Tidy() error {
	return sh.Run("go", "mod", "tidy")
}

func RunLinter() error {
	return sh.Run("golangci-lint", "run")
}

func RunTests() error {
	return sh.Run("go", "test", "./...")
}

func Build() error {
	if err := Tidy(); err != nil {
		return err
	}
	return sh.Run("go", "build", "-o", "bot", "cmd/main.go")
}