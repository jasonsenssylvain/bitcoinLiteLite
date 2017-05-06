package block

import (
	"fmt"
	"testing"

	"time"

	"github.com/jasoncodingnow/bitcoinLiteLite/blockchainuser"
	"github.com/jasoncodingnow/bitcoinLiteLite/tool"
)

func Test_BlockChain(t *testing.T) {
	fmt.Println("start test blockchain")
	b := NewBlockChain()
	newB := b.NewBlock()
	b.Block = newB

	t1 := NewTransaction([]byte(blockchainuser.GetKey().PublicKey), nil, []byte(tool.RandomString(100)))
	t2 := NewTransaction([]byte(blockchainuser.GetKey().PublicKey), nil, []byte(tool.RandomString(100)))
	t3 := NewTransaction([]byte(blockchainuser.GetKey().PublicKey), nil, []byte(tool.RandomString(100)))
	t4 := NewTransaction([]byte(blockchainuser.GetKey().PublicKey), nil, []byte(tool.RandomString(100)))
	t5 := NewTransaction([]byte(blockchainuser.GetKey().PublicKey), nil, []byte(tool.RandomString(100)))
	t6 := NewTransaction([]byte(blockchainuser.GetKey().PublicKey), nil, []byte(tool.RandomString(100)))
	t7 := NewTransaction([]byte(blockchainuser.GetKey().PublicKey), nil, []byte(tool.RandomString(100)))

	prefixMatch := tool.GenerateBytes(1, 0)

	t1.Header.Nonce = t1.GenerateNonce(prefixMatch)
	t1.Signature = t1.Sign(blockchainuser.GetKey().PrivateKey)
	t2.Header.Nonce = t2.GenerateNonce(prefixMatch)
	t2.Signature = t2.Sign(blockchainuser.GetKey().PrivateKey)
	t3.Header.Nonce = t3.GenerateNonce(prefixMatch)
	t3.Signature = t3.Sign(blockchainuser.GetKey().PrivateKey)
	t4.Header.Nonce = t4.GenerateNonce(prefixMatch)
	t4.Signature = t4.Sign(blockchainuser.GetKey().PrivateKey)
	t5.Header.Nonce = t5.GenerateNonce(prefixMatch)
	t5.Signature = t5.Sign(blockchainuser.GetKey().PrivateKey)
	t6.Header.Nonce = t6.GenerateNonce(prefixMatch)
	t6.Signature = t6.Sign(blockchainuser.GetKey().PrivateKey)
	t7.Header.Nonce = t7.GenerateNonce(prefixMatch)
	t7.Signature = t7.Sign(blockchainuser.GetKey().PrivateKey)

	b.Run()
	fmt.Println("bc start run")
	b.TransactionChan <- t1
	b.TransactionChan <- t2
	b.TransactionChan <- t3
	b.TransactionChan <- t4
	b.TransactionChan <- t5
	b.TransactionChan <- t6
	b.TransactionChan <- t7

	// before test, should change block time span to 10s
	time.Sleep(20 * time.Second)

	for _, bc := range *b.BlockSlice {
		fmt.Println(bc.Signature)
	}
}
