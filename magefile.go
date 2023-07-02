//go:build mage

package main

import (
	"github.com/magefile/mage/sh"
)

func Build() error {
	if err := sh.Run("go", "mod", "tidy"); err != nil {
		return err
	}
	if err := sh.Run("swag", "init", "-g", "internal/server/handler/handler.go", "--ot", "yaml"); err != nil {
		return err
	}
	if err := sh.Run("go", "test", "./..."); err != nil {
		return err
	}
	return sh.Run("go", "build", "-o", "multimessenger_bot", "cmd/main.go")
}
