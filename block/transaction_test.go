package block

import (
	"testing"
)

func Test_All(t *testing.T) {
	// key := crypto.GenerateKey()
	// transaction := NewTransaction([]byte(key.PublicKey), nil, []byte(tool.RandomString(100)))

	// // 模拟难度
	// prefixMatch := tool.GenerateBytes(1, 0)
	// transaction.Header.Nonce = transaction.GenerateNonce(prefixMatch)
	// transaction.Signature = transaction.Sign(key.PrivateKey)

	// data, err := transaction.MarshalBinary()
	// if err != nil {
	// 	t.Error(err)
	// }

	// // 从data还原回transaction
	// tr2 := &Transaction{}
	// _, err = tr2.UnmarshalBinary(data)
	// if err != nil {
	// 	t.Error(err)
	// }

	// // 对比是否一致
	// if !reflect.DeepEqual(tr2.Signature, transaction.Signature) ||
	// 	!reflect.DeepEqual(tr2.Payload, transaction.Payload) ||
	// 	!reflect.DeepEqual(tr2.Header.From, transaction.Header.From) ||
	// 	!reflect.DeepEqual(tr2.Header.To, transaction.Header.To) ||
	// 	tr2.Header.Nonce != transaction.Header.Nonce ||
	// 	tr2.Header.Timestamp != transaction.Header.Timestamp ||
	// 	!reflect.DeepEqual(tr2.Header.PayloadHash, transaction.Header.PayloadHash) ||
	// 	tr2.Header.PayloadLen != transaction.Header.PayloadLen {

	// 	fmt.Println(reflect.DeepEqual(tr2.Signature, transaction.Signature))
	// 	fmt.Println(reflect.DeepEqual(tr2.Payload, transaction.Payload))
	// 	fmt.Println(reflect.DeepEqual(tr2.Header.From, transaction.Header.From))
	// 	fmt.Println(reflect.DeepEqual(tr2.Header.To, transaction.Header.To))
	// 	fmt.Println(tr2.Header.Nonce != transaction.Header.Nonce)
	// 	fmt.Println(tr2.Header.Timestamp != transaction.Header.Timestamp)
	// 	fmt.Println(tr2.Header.Timestamp)
	// 	fmt.Println(transaction.Header.Timestamp)

	// 	fmt.Println("tr2.Signature: " + string(tr2.Signature))
	// 	fmt.Println("transaction.Signature: " + string(transaction.Signature))
	// 	t.Error("error")
	// }

	// // 检测交易是否合法
	// verifyResult := transaction.VerifyTransaction(prefixMatch)
	// if !verifyResult {
	// 	t.Error("verify failed")
	// }

}
