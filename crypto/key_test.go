package crypto

import (
	"fmt"
	"testing"
)

func Test_All(t *testing.T) {
	keypair := GenerateKey()
	fmt.Println(string(keypair.PrivateKey))
	fmt.Println(string(keypair.PublicKey))
	message := "fwfwfwf"
	signature, _ := Sign([]byte(message), keypair.PrivateKey)
	fmt.Println(string(signature))
	result, _ := Verify([]byte(message), signature, keypair.PublicKey)
	fmt.Println(result)
}
