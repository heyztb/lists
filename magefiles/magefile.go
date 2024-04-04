package main

import (
	"fmt"
	"os"
	"strings"

	_ "github.com/joho/godotenv/autoload"
	"github.com/magefile/mage/sh"
)

func Unit() error {
	args := []string{
		"test",
		"github.com/heyztb/lists-backend/internal/crypto",
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

func Integration() error {
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

func Templ() error {
	err := sh.RunV("templ", "generate")
	if err != nil {
		fmt.Printf("error generating templates %s", err)
		return err
	}
	return nil
}

func Run() error {
	environ := os.Environ()
	env := make(map[string]string, len(environ))
	for _, v := range environ {
		kv := strings.Split(v, "=")
		if len(kv) == 2 {
			env[kv[0]] = kv[1]
		}
	}

	fmt.Printf("Starting server on %s\n", env["LISTEN_ADDRESS"])
	err := sh.RunWithV(env, "go", "run", "./cmd/backend")
	if err != nil {
		fmt.Printf("error running server %s", err)
		return err
	}

	return nil
}
