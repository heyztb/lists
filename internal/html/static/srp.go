//go:build wasm
// +build wasm

package static

import (
	"crypto"
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

func throwError(err error) {
	js.Global().Call("eval", fmt.Sprintf(`throw new Error("%s")`, err.Error()))
}

func getRegistrationInfo() js.Func {
	return js.FuncOf(func(_ js.Value, args []js.Value) any {
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
	})
}

func setupClient() js.Func {
	return js.FuncOf(func(_ js.Value, args []js.Value) any {
		if len(args) != 4 {
			throwError(fmt.Errorf("setupClient: Invalid number of arguments (expects 3, got %d)", len(args)))
			return nil
		}

		identifier := args[0].String()
		password := args[1].String()
		salt := make([]byte, srp.SaltLength)
		B := make([]byte, 768)
		js.CopyBytesToGo(salt, args[2])
		js.CopyBytesToGo(B, args[3])

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
	})
}

func getClientProof() js.Func {
	return js.FuncOf(func(_ js.Value, args []js.Value) any {
		proof, err := client.ComputeM1()
		if err != nil {
			throwError(fmt.Errorf("proof: failed to generate client proof -- %w", err))
			return nil
		}

		return hex.EncodeToString(proof)
	})
}

func verifyServerProof() js.Func {
	return js.FuncOf(func(_ js.Value, args []js.Value) any {
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
	})
}

func getKey() js.Func {
	return js.FuncOf(func(_ js.Value, args []js.Value) any {
		key, err := client.SessionKey()
		if err != nil {
			throwError(fmt.Errorf("getKey: error getting session key -- %w", err))
			return nil
		}

		return hex.EncodeToString(key)
	})
}

func main() {
	c := make(chan struct{}, 0)
	js.Global().Set("getRegistrationInfo", getRegistrationInfo)
	js.Global().Set("setupClient", setupClient)
	js.Global().Set("getClientProof", getClientProof)
	js.Global().Set("verifyServerProof", verifyServerProof)
	js.Global().Set("getKey", getKey)
	<-c
	getRegistrationInfo().Release()
	setupClient().Release()
	getClientProof().Release()
	verifyServerProof().Release()
	getKey().Release()
}