package block

import "github.com/jasoncodingnow/bitcoinLiteLite/crypto"

const PayloadHashSize = 80
const TransactionSignatureSize = 80

const TransactionHeaderSize = crypto.PublicKeyLen + crypto.PublicKeyLen + PayloadHashSize + 4 + 4 + 4

const MaxInt = int(^uint(0) >> 1)

const BlockSignatureSize = 80
const MerkleRootSize = 80
const BlockHeaderSize = crypto.PublicKeyLen + BlockSignatureSize + MerkleRootSize + 4 + 4

var BlockPowPrefix = 2       // 该参数应该随着难度调整，打包的难度
var TransactionPowPrefix = 1 // 交易的难度，这个目前应该没什么太大用处

const PackageTimespan = 1 //1分钟打一个block
