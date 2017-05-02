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

const (
	SignRLen = 28
	SignSLen = 28
)

type Key struct {
	PrivateKey []byte
	PublicKey  []byte
}

//GenerateKey 生成对应的公钥和私钥
func GenerateKey() *Key {
	key := &Key{}
	pk, _ := ecdsa.GenerateKey(elliptic.P224(), rand.Reader)
	x509encoded, _ := x509.MarshalECPrivateKey(pk)
	key.PrivateKey = pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: x509encoded})

	x509encodedpub, _ := x509.MarshalPKIXPublicKey(&pk.PublicKey)
	key.PublicKey = pem.EncodeToMemory(&pem.Block{Type: "EC PUBLIC KEY", Bytes: x509encodedpub})
	return key
}

//Sign 通过私钥对数据进行签名
func Sign(payload []byte, pk []byte) ([]byte, error) {
	block, _ := pem.Decode(pk)
	x509encoded := block.Bytes
	privateKey, err := x509.ParseECPrivateKey(x509encoded)
	if err != nil {
		return nil, err
	}
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, payload)
	rBytes := tool.FillBytesToFront(r.Bytes(), SignRLen)
	sBytes := tool.FillBytesToFront(s.Bytes(), SignSLen)
	signature := append(rBytes, sBytes...)

	return signature, nil
}

//Verify 验证签名对否
func Verify(payload, signature []byte, publicKey []byte) (bool, error) {
	blockPub, _ := pem.Decode(publicKey)
	x509EncodedPub := blockPub.Bytes
	genericPublicKey, _ := x509.ParsePKIXPublicKey(x509EncodedPub)
	genericPublicKeyA := genericPublicKey.(*ecdsa.PublicKey)

	buf := bytes.NewBuffer(signature)
	rBytes := new(big.Int).SetBytes(buf.Next(SignRLen))
	sBytes := new(big.Int).SetBytes(buf.Next(SignSLen))
	verifystatus := ecdsa.Verify(genericPublicKeyA, payload, rBytes, sBytes)
	return verifystatus, nil
}
