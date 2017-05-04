package crypto

import (
	"fmt"
	"testing"
)

func Test_All(t *testing.T) {

	keypair := GenerateKey()
	fmt.Println(string(keypair.PrivateKey))
	fmt.Println(len(keypair.PrivateKey))
	fmt.Println(string(keypair.PublicKey))
	fmt.Println(len(keypair.PublicKey))
	message := "fwfwwwwwwwfooooooooooooooooooo33333wf"
	signature, _ := Sign([]byte(message), keypair.PrivateKey)
	fmt.Println(string(signature))
	fmt.Println(len(signature))
	result := Verify([]byte(message), signature, keypair.PublicKey)
	fmt.Println(result)

}
