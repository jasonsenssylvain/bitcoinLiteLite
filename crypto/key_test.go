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

	// a := 1

	// buf := bytes.Buffer{}
	// binary.Write(&buf, binary.BigEndian, a)
	// fmt.Println(buf.Bytes())

	// var newB = make([]byte, 4)
	// binary.LittleEndian.PutUint32(newB, uint32(a))
	// fmt.Println(newB)

	// b := int(binary.LittleEndian.Uint32(newB))
	// fmt.Println(b)
}
