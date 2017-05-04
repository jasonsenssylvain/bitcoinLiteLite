package btc

import (
	"errors"
	"time"

	"bytes"

	"encoding/binary"

	"reflect"

	"github.com/jasoncodingnow/bitcoinLiteLite/consensus"
	"github.com/jasoncodingnow/bitcoinLiteLite/crypto"
	"github.com/jasoncodingnow/bitcoinLiteLite/tool"
)

type TransactionHeader struct {
	From        []byte
	To          []byte
	PayloadHash []byte
	PayloadLen  uint32
	Timestamp   uint32
	Nonce       uint32
}

type Transaction struct {
	Header    *TransactionHeader
	Signature []byte
	Payload   []byte
}

type TransactionSlice []Transaction

//NewTransaction create new transaction
func NewTransaction(from, to, payload []byte) *Transaction {
	t := &Transaction{Header: &TransactionHeader{From: from, To: to}, Payload: payload}
	t.Header.Timestamp = uint32(time.Now().Unix())
	t.Header.PayloadHash = tool.SHA256(payload)
	t.Header.PayloadLen = uint32(len(payload))
	return t
}

func (t *Transaction) Hash() []byte {
	headerBytes, _ := t.Header.MarshalBinary()
	return tool.SHA256(headerBytes)
}

func (t *Transaction) Sign(privateKey string) []byte {
	s, _ := crypto.Sign(t.Hash(), privateKey)
	return []byte(s)
}

func (t *Transaction) VerifyTransaction(powPrefix []byte) bool {
	h := t.Hash()
	payloadHash := tool.SHA256(t.Payload)
	return reflect.DeepEqual(payloadHash, t.Header.PayloadHash) && consensus.CheckProofOfWork(powPrefix, h) && crypto.Verify(h, string(t.Signature), string(t.Header.From))
}

//GenerateNonce pow过程，系统必须不断计算，产生符合条件的nonce才能打包
func (t *Transaction) GenerateNonce(powPrefix []byte) uint32 {
	for {
		if consensus.CheckProofOfWork(powPrefix, t.Hash()) {
			break
		}
		t.Header.Nonce++
	}
	return t.Header.Nonce
}

func (t *Transaction) MarshalBinary() ([]byte, error) {
	headerBinary, err := t.Header.MarshalBinary()
	if err != nil {
		return nil, err
	}
	if len(headerBinary) != TransactionHeaderSize {
		return nil, errors.New("header marshal len error")
	}
	return append(append(headerBinary, tool.FillBytesToFront(t.Signature, TransactionSignatureSize)...), t.Payload...), nil
}

func (t *Transaction) UnmarshalBinary(data []byte) ([]byte, error) {
	buf := bytes.NewBuffer(data)

	if len(data) < (TransactionHeaderSize + TransactionSignatureSize) {
		return nil, errors.New("data length error when unmarshal binary to transaction")
	}

	h := &TransactionHeader{}
	if err := h.UnmarshalBinary(buf.Next(TransactionHeaderSize)); err != nil {
		return nil, err
	}

	t.Header = h
	t.Signature = tool.SliceByteWhenEncount(buf.Next(TransactionSignatureSize), 0)
	t.Payload = buf.Next(int(t.Header.PayloadLen))
	return buf.Next(MaxInt), nil
}

func (t *TransactionHeader) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	buf.Write(tool.FillBytesToFront(t.From, crypto.PublicKeyLen))
	buf.Write(tool.FillBytesToFront(t.To, crypto.PublicKeyLen))
	buf.Write(tool.FillBytesToFront(t.PayloadHash, PayloadHashSize))
	binary.Write(buf, binary.LittleEndian, t.PayloadLen)
	binary.Write(buf, binary.LittleEndian, t.Timestamp)
	binary.Write(buf, binary.LittleEndian, t.Nonce)
	return buf.Bytes(), nil
}

func (t *TransactionHeader) UnmarshalBinary(data []byte) error {
	buf := bytes.NewBuffer(data)
	t.From = tool.SliceByteWhenEncount(buf.Next(crypto.PublicKeyLen), 0)
	t.To = tool.SliceByteWhenEncount(buf.Next(crypto.PublicKeyLen), 0)
	t.PayloadHash = tool.SliceByteWhenEncount(buf.Next(PayloadHashSize), 0)
	binary.Read(bytes.NewBuffer(buf.Next(4)), binary.LittleEndian, &t.PayloadLen)
	binary.Read(bytes.NewBuffer(buf.Next(4)), binary.LittleEndian, &t.Timestamp)
	binary.Read(bytes.NewBuffer(buf.Next(4)), binary.LittleEndian, &t.Nonce)
	return nil
}

func (t *TransactionSlice) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	for _, tr := range *t {
		b, err := tr.MarshalBinary()
		if err != nil {
			return nil, err
		}
		buf.Write(b)
	}
	return buf.Bytes(), nil
}

func (t *TransactionSlice) UnmarshalBinary(data []byte) error {
	d := data

	for len(d) > TransactionHeaderSize+TransactionSignatureSize {
		tr := &Transaction{}
		remain, err := tr.UnmarshalBinary(d)
		if err != nil {
			return err
		}
		*t = append(*t, *tr)
		d = remain
	}
	return nil
}

func (t TransactionSlice) Exists(newTr *Transaction) bool {
	for _, tr := range t {
		if reflect.DeepEqual(tr.Signature, newTr.Signature) {
			return true
		}
	}
	return false
}

func (t TransactionSlice) AddTransaction(newTr *Transaction) TransactionSlice {
	for i, tr := range t {
		if tr.Header.Timestamp >= newTr.Header.Timestamp {
			return append(append(t[:i], *newTr), t[i:]...)
		}
	}
	return append(t, *newTr)
}
