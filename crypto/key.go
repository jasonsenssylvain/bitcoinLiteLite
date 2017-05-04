package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"math/big"

	"bytes"

	"github.com/jasoncodingnow/bitcoinLiteLite/tool"
)

type Key struct {
	PrivateKey string
	PublicKey  string
}

var b58 = NewBitcoinBase58()

//GenerateKey 生成对应的公钥和私钥
func GenerateKey() *Key {
	key := &Key{}
	pk, _ := ecdsa.GenerateKey(elliptic.P224(), rand.Reader)
	x509encoded, _ := x509.MarshalECPrivateKey(pk)
	key.PrivateKey, _ = b58.EncodeToString(pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: x509encoded}))

	x509encodedpub, _ := x509.MarshalPKIXPublicKey(&pk.PublicKey)
	key.PublicKey, _ = b58.EncodeToString(pem.EncodeToMemory(&pem.Block{Type: "EC PUBLIC KEY", Bytes: x509encodedpub}))
	return key
}

//Sign 通过私钥对数据进行签名
func Sign(payload []byte, privateKey string) (string, error) {
	pk, _ := b58.DecodeString(privateKey)
	block, _ := pem.Decode(pk)
	x509encoded := block.Bytes
	realPrivateKey, err := x509.ParseECPrivateKey(x509encoded)
	if err != nil {
		return "", err
	}
	r, s, err := ecdsa.Sign(rand.Reader, realPrivateKey, payload)
	rBytes := tool.FillBytesToFront(r.Bytes(), SignRLen)
	sBytes := tool.FillBytesToFront(s.Bytes(), SignSLen)
	signature := append(rBytes, sBytes...)

	return b58.EncodeToString(signature)
}

//Verify 验证签名对否
func Verify(payload []byte, signature string, publicKey string) bool {
	pk, _ := b58.DecodeString(publicKey)

	blockPub, _ := pem.Decode(pk)
	x509EncodedPub := blockPub.Bytes
	genericPublicKey, _ := x509.ParsePKIXPublicKey(x509EncodedPub)
	genericPublicKeyA := genericPublicKey.(*ecdsa.PublicKey)

	sign, _ := b58.DecodeString(signature)
	buf := bytes.NewBuffer(sign)
	rBytes := new(big.Int).SetBytes(buf.Next(SignRLen))
	sBytes := new(big.Int).SetBytes(buf.Next(SignSLen))
	verifystatus := ecdsa.Verify(genericPublicKeyA, payload, rBytes, sBytes)
	return verifystatus
}
