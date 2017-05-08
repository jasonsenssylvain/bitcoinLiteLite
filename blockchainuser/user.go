package blockchainuser

import (
	"github.com/jasoncodingnow/bitcoinLiteLite/crypto"
)

var key *crypto.Key

//GetKey 获取当前用户私钥公钥
func GetKey() *crypto.Key {
	if key == nil || key.PrivateKey == "" {
		key = crypto.GenerateKey()
	}
	return key
}

func SetKey(k *crypto.Key) {
	key = k
}
