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

	blockBytes, err := b.MarshalBinary()
	fmt.Println("len of blockBytes")
	fmt.Println(len(blockBytes))
	if err != nil {
		fmt.Println(err)
	}
	newBlock := NewBlock(nil)
	newBlock.UnmarshalBinary(blockBytes)

	//对比
	fmt.Println("block Merkle Root ")
	fmt.Println(b.Header.MerkleRoot)
	fmt.Println("newBlock Merkle Root ")
	fmt.Println(newBlock.Header.MerkleRoot)

	fmt.Println("block Origin ")
	fmt.Println(len(b.Header.Origin))
	fmt.Println(b.Header.Origin)
	fmt.Println("newBlock Origin ")
	fmt.Println(newBlock.Header.Origin)

	fmt.Println("block PrevBlock ")
	fmt.Println(len(b.Header.PrevBlock))
	fmt.Println(b.Header.PrevBlock)
	fmt.Println("newBlock PrevBlock ")
	fmt.Println(newBlock.Header.PrevBlock)
}
