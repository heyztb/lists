// small utility to generate a new Paseto V4 Asymmetric secret key -- export PASETO_KEY=`go run cmd/generate-paseto-key/main.go`
package main

import (
	"fmt"

	"aidanwoods.dev/go-paseto"
)

func main() {
	key := paseto.NewV4AsymmetricSecretKey()
	fmt.Println(key.ExportHex())
}
