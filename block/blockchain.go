package block

import (
	"fmt"
	"strconv"
	"time"

	"reflect"

	"github.com/jasoncodingnow/bitcoinLiteLite/blockchainuser"
	"github.com/jasoncodingnow/bitcoinLiteLite/consensus"
	"github.com/jasoncodingnow/bitcoinLiteLite/message"
	"github.com/jasoncodingnow/bitcoinLiteLite/network"
	"github.com/jasoncodingnow/bitcoinLiteLite/tool"
)

type BlockChan chan Block
type TransactionChan chan Transaction

type BlockChain struct {
	Block      *Block // 当前的Block，属于未广播未确认状态
	BlockSlice *BlockSlice

	BlockChan         BlockChan        // 接收Block
	TransactionChan   TransactionChan  // 接收Transaction，验证，加入到Block
	RemainTransactins TransactionSlice // 正在打包，或者正在验证区块，如果整个时候有新的Transaction进入，则暂时放入这里
}

var isPackaging = false // 是否当前正在打包

func NewBlockChain() *BlockChain {
	b := &BlockChain{}
	b.Block = NewBlock(nil)
	b.BlockSlice = &BlockSlice{}
	b.BlockChan = make(BlockChan)
	b.TransactionChan = make(TransactionChan)
	b.RemainTransactins = TransactionSlice{}
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
	l := len(*bc.BlockSlice)
	if l != 0 {
		b.Header.PrevBlock = (*bc.BlockSlice)[l-1].Hash()
	}
	newBlockSlice := append(*bc.BlockSlice, *b)
	bc.BlockSlice = &newBlockSlice
}

func (bc *BlockChain) Run() {
	go func() {
		newBlockChan := bc.GenerateBlock()
		bc.newTicker(newBlockChan)
		for {
			select {
			case tr := <-bc.TransactionChan:
				fmt.Println("[INFO] receive new transaction, Signature is : ")
				fmt.Println(tr.Signature)

				if bc.Block.Transactions.Exists(&tr) {
					continue
				}
				if !tr.VerifyTransaction(tool.GenerateBytes(TransactionPowPrefix, 0)) {
					fmt.Println("not valid transaction")
					continue
				}
				if isPackaging {
					// 如果正在打包
					// 如果直接加入的时候，有可能导致检验不通过。这里处理只是减少概率而已，没治本
					bc.RemainTransactins = append(bc.RemainTransactins, tr)
					bc.broadCastTransaction(&tr)
					continue
				}
				bc.Block.AddTransaction(&tr)
				bc.broadCastTransaction(&tr)
				if bc.checkNeedToPackageBlock() {
					bc.addTrToBlock(bc.Block, &bc.RemainTransactins)
					newBlockChan <- *(bc.Block)
				}
			case b := <-bc.BlockChan:
				fmt.Println("receive new block")
				if bc.BlockSlice.Exists(&b) {
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
					fmt.Println("block success validate")
					bc.AppendBlock(&b)

					diffTransactions := bc.getDiffTransactions(bc.Block.Transactions, b.Transactions, &bc.RemainTransactins)

					bc.broadCastBlock(&b)

					bc.Block = bc.NewBlock()
					bc.addTrToBlock(bc.Block, &diffTransactions)
					bc.RemainTransactins = diffTransactions

					fmt.Println("[INFO] after Block package, the len of bc.Block.Transactions is " + strconv.Itoa(len(*bc.Block.Transactions)))
					fmt.Println("[INFO] after Block package, the len of bc.RemainTransactions is " + strconv.Itoa(len(bc.RemainTransactins)))
					if isPackaging {
						isPackaging = false
					}
				}
			}
		}
	}()
}

// 获取两个Transaction之间的不同
func (bc *BlockChain) getDiffTransactions(tr1, tr2, tr3 *TransactionSlice) TransactionSlice {
	result := TransactionSlice{}
	tr1Map := make(map[string]Transaction)
	for _, t := range *tr1 {
		tr1Map[string(t.Signature)] = t
	}

	for _, t := range *tr2 {
		if v, ok := tr1Map[string(t.Signature)]; !ok {
			result = append(result, v)
		}
	}

	for _, t := range *tr3 {
		if v, ok := tr1Map[string(t.Signature)]; !ok {
			result = append(result, v)
		}
	}

	return result
}

// 把trs里面的交易，取出一部分放到Block里
func (bc *BlockChain) addTrToBlock(b *Block, trs *TransactionSlice) {

	tr1Map := make(map[string]Transaction)
	for _, t := range *b.Transactions {
		tr1Map[string(t.Signature)] = t
	}

	result := TransactionSlice{}

	for _, t := range *trs {
		if v, ok := tr1Map[string(t.Signature)]; !ok {
			result = append(result, v)
		}
	}

	trl := len(*b.Transactions)
	noEnough := false
	if trl < 5 {
		for i := 0; i < (5 - trl); i++ {
			if i < len(result) {
				t := append(*b.Transactions, result[i])
				b.Transactions = &t
			} else {
				noEnough = true
			}
		}
	}
	if !noEnough {
		result = result[(5 - trl):]
		trs = &result
	} else {
		trs = &TransactionSlice{}
	}
}

// 当前需要打包有3个条件：
// 1. 当前不是正在打包阶段
// 2. 当前交易数量达到5，但是等待时间无所谓
// // 3. 当前交易数量无所谓，但是距离上一个打包已经十分钟了
func (bc *BlockChain) checkNeedToPackageBlock() bool {
	if isPackaging {
		return false
	}
	if (len(*bc.Block.Transactions) + len(bc.RemainTransactins)) >= 5 {
		return true
	}
	return false
}

// 3. 当前交易数量无所谓，但是距离上一个打包已经十分钟了
func (bc *BlockChain) newTicker(newBlockChan chan Block) {
	go func() {
		timer := time.NewTicker(PackageTimespan * time.Minute)
		// timer := time.NewTicker(10 * time.Second) // test
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
						if (now - lastBlock.Header.Timestamp) >= PackageTimespan*60 {
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
		for true {
			newBlock := <-b
			fmt.Println("new block! start pow")
			calNewBlockFinish := false
			for !calNewBlockFinish {
				isPackaging = true

				// 初始化Block Header
				if newBlock.Header.Origin == nil {
					newBlock.Header.Origin = []byte(blockchainuser.GetKey().PublicKey)
				}
				newBlock.Header.MerkleRoot = newBlock.GenrateMerkleRoot()
				newBlock.Header.Timestamp = uint32(time.Now().Unix())
				newBlock.Header.Nonce = 0

				powSuccess := false
				// POW
				// 如果有新的Block进入，并且正好在isPackaging，则从外界打断该过程
				//TODO: 这里应该有更好的机制，需要完善POW过程中，有新的Block验证完成，如何处理
				for isPackaging {
					if consensus.CheckProofOfWork(tool.GenerateBytes(BlockPowPrefix, 0), newBlock.Hash()) {
						newBlock.Signature = newBlock.Sign(blockchainuser.GetKey().PrivateKey)
						bc.BlockChan <- newBlock // 自己产生的区块也发给同一个channel进行验证和写入
						calNewBlockFinish = true
						fmt.Println("new block packaged!")
						powSuccess = true
						break
					}
					newBlock.Header.Nonce++
				}

				if !powSuccess {
					// 说明pow被打断了
					// 判断是否有足够的Transaction，有的话开始打包
					trl := len(*(bc.Block.Transactions))
					if trl >= 5 {
						b <- *(bc.Block)
					}
				}
			}
			isPackaging = false
		}
	}()
	return b
}

// 程序运行到这里，POW完成，需要广播
func (bc *BlockChain) broadCastBlock(b *Block) {
	m, _ := message.NewMessage(message.MessageTypeSendBlock)
	m.Data, _ = b.MarshalBinary()
	network.NetworkInstant.BroadCastChan <- *m
}

// 广播Transaction
func (bc *BlockChain) broadCastTransaction(t *Transaction) {
	m, _ := message.NewMessage(message.MessageTypeSendTransaction)
	m.Data, _ = t.MarshalBinary()
	network.NetworkInstant.BroadCastChan <- *m
}
