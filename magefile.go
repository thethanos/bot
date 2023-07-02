//go:build mage

package main

import (
	"github.com/magefile/mage/sh"
)

func Tidy() error {
	return sh.Run("go", "mod", "tidy")
}

func GenDoc() error {
	return sh.Run("swag", "init", "-g", "internal/server/handler/handler.go", "--ot", "yaml")
}

func RunTests() error {
	return sh.Run("go", "test", "./...")
}

func Build() error {
	if err := Tidy(); err != nil {
		return err
	}
	if err := GenDoc(); err != nil {
		return err
	}
	if err := RunTests(); err != nil {
		return err
	}
	return sh.Run("go", "build", "-o", "multimessenger_bot", "cmd/main.go")
}
