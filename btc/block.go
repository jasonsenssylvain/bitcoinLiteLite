package btc

import (
	"bytes"
	"encoding/binary"
	"time"

	"reflect"

	"github.com/jasoncodingnow/bitcoinLiteLite/consensus"
	"github.com/jasoncodingnow/bitcoinLiteLite/crypto"
	"github.com/jasoncodingnow/bitcoinLiteLite/tool"
)

type Block struct {
	Header       *BlockHeader
	Signature    []byte
	Transactions *TransactionSlice
}

type BlockHeader struct {
	Origin     []byte // public key or address
	PrevBlock  []byte
	MerkleRoot []byte
	Timestamp  uint32
	Nonce      uint32
}

type BlockSlice []Block

func NewBlock(prevBlock []byte) *Block {
	b := &Block{}
	b.Header = &BlockHeader{PrevBlock: prevBlock}
	b.Header.Timestamp = uint32(time.Now().Unix())
	return b
}

func (b *Block) AddTransaction(t *Transaction) {
	newSlice := b.Transactions.AddTransaction(t)
	b.Transactions = &newSlice
}

func (b *Block) Hash() []byte {
	headerBytes := b.Header.MarshalBinary()
	return tool.SHA256(headerBytes)
}

func (b *Block) GenerateNonce(powPrefix []byte) uint32 {
	for {
		if consensus.CheckProofOfWork(powPrefix, b.Hash()) {
			break
		}
		b.Header.Nonce++
	}
	return b.Header.Nonce
}

func (b *Block) VerifyBlock(powPrefix []byte) bool {
	h := b.Hash()
	m := b.GenrateMerkleRoot()

	return reflect.DeepEqual(m, b.Header.MerkleRoot) && consensus.CheckProofOfWork(powPrefix, h) && crypto.Verify(h, string(b.Signature), string(b.Header.Origin))
}

func (b *Block) Sign(privateKey string) []byte {
	sign, _ := crypto.Sign(b.Hash(), privateKey)
	return []byte(sign)
}

func (b *Block) GenrateMerkleRoot() []byte {
	l := len(*b.Transactions)
	tree := make([][]byte, l)
	for _, t := range *b.Transactions {
		tree = append(tree, t.Hash())
	}
	merkleRoot := b.generateMerkleRoot(tree)
	return merkleRoot
}

// 通过递归生成MerkleTree的最终hash
func (b *Block) generateMerkleRoot(tree [][]byte) []byte {
	l := len(tree)
	if l == 0 {
		return nil
	}
	if l == 1 {
		return tree[0]
	}
	lastTreeNode := []byte{}
	half := 0
	if l%2 == 1 {
		// 奇数，把最后一个拿出来
		lastTreeNode = tree[l-1]
		half = (l - 1) / 2
	} else {
		half = l / 2
	}

	newTree := make([][]byte, 0)
	for i := 0; i < half; i++ {
		prevNode, nextNode := tree[i*2], tree[i*2+1]
		hash := tool.SHA256(append(prevNode, nextNode...))
		newTree = append(newTree, hash)
	}

	if len(lastTreeNode) == 0 {
		return b.generateMerkleRoot(newTree)
	}
	newTree = append(newTree, lastTreeNode)
	return b.generateMerkleRoot(newTree)
}

func (b *Block) MarshalBinary() ([]byte, error) {
	binary := b.Header.MarshalBinary()
	signature := tool.FillBytesToFront(b.Signature, BlockSignatureSize)
	transactionBytes, err := b.Transactions.MarshalBinary()
	if err != nil {
		return nil, err
	}
	return append(append(binary, signature...), transactionBytes...), nil
}

func (b *Block) UnmarshalBinary(data []byte) error {
	buf := bytes.NewBuffer(data)

	err := b.Header.UnmarshalBinary(buf.Next(BlockHeaderSize))
	if err != nil {
		return err
	}

	b.Signature = tool.SliceByteWhenEncount(buf.Next(BlockSignatureSize), 0)
	err = b.Transactions.UnmarshalBinary(buf.Next(MaxInt))
	if err != nil {
		return err
	}
	return nil
}

func (h *BlockHeader) MarshalBinary() []byte {
	buf := &bytes.Buffer{}

	buf.Write(tool.FillBytesToFront(h.Origin, crypto.PublicKeyLen))
	buf.Write(tool.FillBytesToFront(h.PrevBlock, BlockSignatureSize))
	buf.Write(tool.FillBytesToFront(h.MerkleRoot, MerkleRootSize))
	binary.Write(buf, binary.LittleEndian, h.Timestamp)
	binary.Write(buf, binary.LittleEndian, h.Nonce)
	return buf.Bytes()
}

func (h *BlockHeader) UnmarshalBinary(data []byte) error {
	buf := bytes.NewBuffer(data)
	h.Origin = tool.SliceByteWhenEncount(buf.Next(crypto.PublicKeyLen), 0)
	h.Origin = tool.SliceByteWhenEncount(buf.Next(BlockSignatureSize), 0)
	h.Origin = tool.SliceByteWhenEncount(buf.Next(MerkleRootSize), 0)
	binary.Read(bytes.NewBuffer(buf.Next(4)), binary.LittleEndian, &h.Timestamp)
	binary.Read(bytes.NewBuffer(buf.Next(4)), binary.LittleEndian, &h.Nonce)
	return nil
}

func (slice *BlockSlice) Exists(b *Block) bool {
	for i := len(*slice) - 1; i >= 0; i-- {
		currB := (*slice)[i]
		if reflect.DeepEqual(currB.Signature, b.Signature) {
			return true
		}
	}
	return false
}
