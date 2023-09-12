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

func GenDoc() error {
	return sh.Run("swag", "init", "-g", "internal/server/handler/handler.go", "--ot", "yaml", "-o", "docs")
}

func RunTests() error {
	return sh.Run("go", "test", "./...")
}

func BuildBot() error {
	if err := Tidy(); err != nil {
		return err
	}
	return sh.Run("go", "build", "-o", "bot", "cmd/bot/main.go")
}

func BuildServer() error {
	if err := Tidy(); err != nil {
		return err
	}
	if err := GenDoc(); err != nil {
		return err
	}
	return sh.Run("go", "build", "-o", "server", "cmd/server/main.go")
}
