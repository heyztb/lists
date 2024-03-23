package main

import (
	"fmt"

	"github.com/magefile/mage/sh"
)

func RunUnitTests() error {
	args := []string{
		"test",
		"github.com/heyztb/lists-backend/internal/crypto",
		"github.com/heyztb/lists-backend/internal/models",
		"github.com/heyztb/lists-backend/internal/paseto",
		"-cover",
		"-v",
	}
	err := sh.RunV("go", args...)
	if err != nil {
		fmt.Printf("error running tests %s", err)
		return err
	}
	return nil
}

func RunIntegrationTest() error {
	args := []string{
		"test",
		"github.com/heyztb/lists-backend/internal/server",
		"-v",
	}

	err := sh.RunV("go", args...)
	if err != nil {
		fmt.Printf("error running integration test %s", err)
		return err
	}
	return nil
}