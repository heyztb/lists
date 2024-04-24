//go:build wasm
// +build wasm

package main

import (
	"crypto"
	_ "crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"

	"syscall/js"

	"code.posterity.life/srp/v2"
	"golang.org/x/crypto/argon2"
)

var client *srp.Client
var params = &srp.Params{
	Name:  "DH15-SHA256-Argon2",
	Group: srp.RFC5054Group3072,
	Hash:  crypto.SHA256,
	KDF: func(username string, password string, salt []byte) ([]byte, error) {
		p := []byte(username + ":" + password)
		key := argon2.IDKey(p, salt, 1, 64*1024, 4, 32)
		return key, nil
	},
}

func main() {
	c := make(chan struct{}, 0)
	js.Global().Set("getRegistrationInfo", js.FuncOf(getRegistrationInfo))
	js.Global().Set("setupClient", js.FuncOf(setupClient))
	js.Global().Set("getClientProof", js.FuncOf(getClientProof))
	js.Global().Set("verifyServerProof", js.FuncOf(verifyServerProof))
	js.Global().Set("getKey", js.FuncOf(getKey))
	<-c
}

func throwError(err error) {
	js.Global().Call("eval", fmt.Sprintf(`throw new Error("%s")`, err.Error()))
}

func recoverError() {
	if err := recover(); err != nil {
		if jsErr, ok := err.(js.Error); ok {
			// Handle the JavaScript error
			fmt.Println("JavaScript Error:", jsErr.Error())
		} else {
			// Handle other types of errors
			fmt.Println("Unexpected Error:", err)
		}
	}
}

func getRegistrationInfo(_ js.Value, args []js.Value) any {
	defer recoverError()

	if len(args) != 2 {
		throwError(errors.New("register: Invalid number of arguments"))
		return nil
	}

	identifier := args[0].String()
	password := args[1].String()

	triplet, err := srp.ComputeVerifier(params, identifier, password, srp.NewSalt())
	if err != nil {
		throwError(fmt.Errorf("register: error computing verifier -- %w", err))
		return nil
	}

	return map[string]any{
		"salt":     hex.EncodeToString(triplet.Salt()),
		"verifier": hex.EncodeToString(triplet.Verifier()),
	}
}

func setupClient(_ js.Value, args []js.Value) any {
	defer recoverError()

	if len(args) != 4 {
		throwError(fmt.Errorf("setupClient: Invalid number of arguments (expects 3, got %d)", len(args)))
		return nil
	}

	identifier := args[0].String()
	password := args[1].String()
	salt := make([]byte, srp.SaltLength)
	B := make([]byte, 384)
	_ = js.CopyBytesToGo(salt, args[2])
	_ = js.CopyBytesToGo(B, args[3])

	var err error
	client, err = srp.NewClient(params, identifier, password, salt)
	if err != nil {
		throwError(fmt.Errorf("setupClient: failed to create srp client -- %w", err))
	}

	err = client.SetB(B)
	if err != nil {
		throwError(fmt.Errorf("setupClient: failed to set server public key -- %w", err))
		return nil
	}

	return hex.EncodeToString(client.A())
}

func getClientProof(_ js.Value, args []js.Value) any {
	defer recoverError()

	proof, err := client.ComputeM1()
	if err != nil {
		throwError(fmt.Errorf("proof: failed to generate client proof -- %w", err))
		return nil
	}

	return hex.EncodeToString(proof)
}

func verifyServerProof(_ js.Value, args []js.Value) any {
	defer recoverError()

	if len(args) != 1 {
		throwError(fmt.Errorf("verifyServerProof: Invalid number of arguments (expects 1, got %d)", len(args)))
		return nil
	}

	serverProof := make([]byte, 32)
	js.CopyBytesToGo(serverProof, args[0])

	valid, err := client.CheckM2(serverProof)
	if err != nil {
		throwError(fmt.Errorf("verifyServerProof: error checking server proof -- %w", err))
		return nil
	}

	return valid
}

func getKey(_ js.Value, args []js.Value) any {
	defer recoverError()
	key, err := client.SessionKey()
	if err != nil {
		throwError(fmt.Errorf("getKey: error getting session key -- %w", err))
		return nil
	}

	return hex.EncodeToString(key)
}
