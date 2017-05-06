package block

import (
	"fmt"
	"testing"

	"github.com/jasoncodingnow/bitcoinLiteLite/crypto"
	"github.com/jasoncodingnow/bitcoinLiteLite/tool"
)

func Test_Block(t *testing.T) {
	fmt.Println("Test_Block start")
	key := crypto.GenerateKey()
	t1 := NewTransaction([]byte(key.PublicKey), nil, []byte(tool.RandomString(100)))
	t2 := NewTransaction([]byte(key.PublicKey), nil, []byte(tool.RandomString(100)))
	t3 := NewTransaction([]byte(key.PublicKey), nil, []byte(tool.RandomString(100)))
	t4 := NewTransaction([]byte(key.PublicKey), nil, []byte(tool.RandomString(100)))
	t5 := NewTransaction([]byte(key.PublicKey), nil, []byte(tool.RandomString(100)))
	t6 := NewTransaction([]byte(key.PublicKey), nil, []byte(tool.RandomString(100)))
	t7 := NewTransaction([]byte(key.PublicKey), nil, []byte(tool.RandomString(100)))

	b := NewBlock(nil)
	b.Header.Origin = []byte(key.PublicKey)
	b.Transactions = &TransactionSlice{*t1, *t2, *t3, *t4, *t5, *t6, *t7}

	prefixMatch := tool.GenerateBytes(1, 0)

	b.Header.MerkleRoot = b.GenrateMerkleRoot()
	b.GenerateNonce(prefixMatch)
	b.Signature = b.Sign(key.PrivateKey)

	verifyResult := b.VerifyBlock(prefixMatch)
	fmt.Println(verifyResult)
}
