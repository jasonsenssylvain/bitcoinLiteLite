package block

import (
	"fmt"
	"time"

	"reflect"

	"github.com/jasoncodingnow/bitcoinLiteLite/blockchainuser"
	"github.com/jasoncodingnow/bitcoinLiteLite/consensus"
	"github.com/jasoncodingnow/bitcoinLiteLite/tool"
)

type BlockChan chan *Block
type TransactionChan chan *Transaction

type BlockChain struct {
	Block      *Block // 当前的Block，属于未广播未确认状态
	BlockSlice *BlockSlice

	BlockChan       BlockChan       // 接收Block
	TransactionChan TransactionChan // 接收Transaction，验证，加入到Block
	// RemainTransactins TransactionSlice // 正在打包，或者正在验证区块，如果整个时候有新的Transaction进入，则暂时放入这里	// 不需要
}

var isPackaging = false // 是否当前正在打包

func NewBlockChain() *BlockChain {
	b := &BlockChain{}
	b.Block = NewBlock(nil)
	b.BlockChan = make(BlockChan)
	b.TransactionChan = make(TransactionChan)
	return b
}

//NewBlock 产生新的区块
func (b *BlockChain) NewBlock() *Block {
	prevBlockHash := []byte{}
	if *b.BlockSlice == nil || len(*b.BlockSlice) == 0 {
		prevBlockHash = nil
	} else {
		prevBlock := (*b.BlockSlice)[len(*b.BlockSlice)-1]
		prevBlockHash = prevBlock.Hash()
	}

	newB := NewBlock(prevBlockHash)
	newB.Header.Origin = []byte(blockchainuser.GetKey().PublicKey)
	return newB
}

//AppendBlock 并且设置Hash
func (bc *BlockChain) AppendBlock(b *Block) {
	b.Header.PrevBlock = (*bc.BlockSlice)[len(*bc.BlockSlice)-1].Hash()
	newBlockSlice := append(*bc.BlockSlice, *b)
	bc.BlockSlice = &newBlockSlice
}

func (bc *BlockChain) Run() {
	newBlockChan := bc.GenerateBlock()
	bc.newTicker(newBlockChan)
	for {
		select {
		case tr := <-bc.TransactionChan:
			if bc.Block.Transactions.Exists(tr) {
				continue
			}
			if !tr.VerifyTransaction(tool.GenerateBytes(TransactionPowPrefix, 0)) {
				fmt.Println("not valid transaction")
				continue
			}
			// if bc.isPackaging {
			// 	// 如果正在打包

			// }
			bc.Block.AddTransaction(tr)
			if bc.checkNeedToPackageBlock() {
				newBlockChan <- *(bc.Block)
			}
		case b := <-bc.BlockChan:
			if bc.BlockSlice.Exists(b) {
				fmt.Println("block exists")
				continue
			}
			if !b.VerifyBlock(tool.GenerateBytes(BlockPowPrefix, 0)) {
				fmt.Println("block not validate")
				continue
			}
			if reflect.DeepEqual(bc.Block.Hash(), b.Header.PrevBlock) {
				fmt.Println("block not match")
			} else {
				bc.AppendBlock(b)

				diffTransactions := bc.getDiffTransactions(bc.Block.Transactions, b.Transactions)

				bc.broadCastBlock(b)

				bc.Block = bc.NewBlock()
				bc.Block.Transactions = &diffTransactions

			}
		}
	}
}

// 获取两个Transaction之间的不同
func (bc *BlockChain) getDiffTransactions(tr1, tr2 *TransactionSlice) TransactionSlice {
	result := TransactionSlice{}
	var tr1Map map[string]Transaction
	for _, t := range *tr1 {
		tr1Map[string(t.Signature)] = t
	}

	for _, t := range *tr2 {
		if v, ok := tr1Map[string(t.Signature)]; !ok {
			result = append(result, v)
		}
	}

	return result
}

// 当前需要打包有3个条件：
// 1. 当前不是正在打包阶段
// 2. 当前交易数量达到5，但是等待时间无所谓
// // 3. 当前交易数量无所谓，但是距离上一个打包已经十分钟了
func (bc *BlockChain) checkNeedToPackageBlock() bool {
	if isPackaging {
		return false
	}
	if len(*bc.Block.Transactions) >= 5 {
		return true
	}
	return false
}

// 3. 当前交易数量无所谓，但是距离上一个打包已经十分钟了
func (bc *BlockChain) newTicker(newBlockChan chan Block) {
	go func() {
		timer := time.NewTicker(10 * time.Minute)
		for {
			select {
			case <-timer.C:
				l := len(*(bc.BlockSlice))
				if len(*(bc.Block.Transactions)) > 0 {
					if l == 0 {
						newBlockChan <- *(bc.Block)
					} else {
						lastBlock := (*(bc.BlockSlice))[l-1]
						now := uint32(time.Now().Unix())
						if (now - lastBlock.Header.Timestamp) > 10*60000 {
							newBlockChan <- *(bc.Block)
						}
					}
				}
			}
		}
	}()
}

//GenerateBlock 产生新的block，即 打包
//TODO: 需要有详细的打包规则，比如打包是交易达到一定数量开始打包，后续的交易进入下个区块？
func (bc *BlockChain) GenerateBlock() chan Block {
	isPackaging = false
	b := make(chan Block)
	go func() {
		newBlock := <-b
		fmt.Println("new block! start pow")
		calNewBlockFinish := false
		for !calNewBlockFinish {
			isPackaging = true
			// 初始化Block Header
			newBlock.Header.MerkleRoot = newBlock.GenrateMerkleRoot()
			newBlock.Header.Timestamp = uint32(time.Now().Unix())
			newBlock.Header.Nonce = 0

			// POW
			for true {
				if consensus.CheckProofOfWork(tool.GenerateBytes(BlockPowPrefix, 0), newBlock.Hash()) {
					newBlock.Signature = newBlock.Sign(blockchainuser.GetKey().PrivateKey)
					bc.BlockChan <- &newBlock // 自己产生的区块也发给同一个channel进行验证和写入
					calNewBlockFinish = true
					fmt.Println("new block packaged!")
					break
				}
				newBlock.Header.Nonce++
			}

			// 程序运行到这里，POW完成，需要广播
			bc.broadCastBlock(&newBlock)
		}
		isPackaging = false
	}()
	return b
}

// 程序运行到这里，POW完成，需要广播
func (bc *BlockChain) broadCastBlock(b *Block) {

}
