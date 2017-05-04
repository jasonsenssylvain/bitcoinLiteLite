package btc

type Block struct {
	Header *BlockHeader
}

type BlockHeader struct {
	Origin     []byte
	PrevBlock  []byte
	MerkleRoot []byte
	Timestamp  uint32
	Nonce      uint32
}
