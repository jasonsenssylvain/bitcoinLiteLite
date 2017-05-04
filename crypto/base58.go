package crypto

import (
	"bytes"
	"errors"
	"math/big"
)

const (
	BitcoinBase58Chars = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
	Base               = 58
)

var (
	ErrInvalidLengthBytes = errors.New("invalid length bytes")
	ErrInvalidChar        = errors.New("invalid char")
)

type Base58 struct {
	chars      [58]byte
	charIdxMap map[byte]int64
}

func NewBase58(charsStr string) (*Base58, error) {
	b58 := &Base58{}

	if err := b58.initChars(charsStr); err != nil {
		return nil, err
	}

	if err := b58.initCharIdxMap(charsStr); err != nil {
		return nil, err
	}

	return b58, nil
}

func NewBitcoinBase58() *Base58 {
	b58, _ := NewBase58(BitcoinBase58Chars)

	return b58
}

func (b58 *Base58) initChars(charsStr string) error {
	if len(charsStr) != 58 {
		return ErrInvalidLengthBytes
	}

	chars := []byte(charsStr)
	copy(b58.chars[:], chars[:])

	return nil
}

func (b58 *Base58) initCharIdxMap(charsStr string) error {
	if len(charsStr) != 58 {
		return ErrInvalidLengthBytes
	}

	b58.charIdxMap = map[byte]int64{}
	for i, b := range []byte(charsStr) {
		b58.charIdxMap[b] = int64(i)
	}

	return nil
}

func (b58 *Base58) EncodeToString(b []byte) (string, error) {
	n := &big.Int{}
	n.SetBytes(b)

	buf := &bytes.Buffer{}
	for _, c := range b {
		if c == 0x00 {
			if err := buf.WriteByte(b58.chars[0]); err != nil {
				return "", err
			}
		} else {
			break
		}
	}

	zero := big.NewInt(0)
	div := big.NewInt(Base)
	mod := &big.Int{}

	tmpBuf := &bytes.Buffer{}
	for {
		if n.Cmp(zero) == 0 {
			tmpBytes := tmpBuf.Bytes()
			for i := len(tmpBytes) - 1; i >= 0; i-- {
				buf.WriteByte(tmpBytes[i])
			}
			return buf.String(), nil
		}

		n.DivMod(n, div, mod)
		if err := tmpBuf.WriteByte(b58.chars[mod.Int64()]); err != nil {
			return "", err
		}
	}
}

func (b58 *Base58) DecodeString(s string) ([]byte, error) {
	b := []byte(s)

	startIdx := 0
	buf := &bytes.Buffer{}
	for i, c := range b {
		if c == b58.chars[0] {
			if err := buf.WriteByte(0x00); err != nil {
				return nil, err
			}
		} else {
			startIdx = i
			break
		}
	}

	n := big.NewInt(0)
	div := big.NewInt(Base)

	for _, c := range b[startIdx:] {
		charIdx, ok := b58.charIdxMap[c]
		if !ok {
			return nil, ErrInvalidChar
		}

		n.Add(n.Mul(n, div), big.NewInt(charIdx))
	}

	buf.Write(n.Bytes())

	return buf.Bytes(), nil
}
