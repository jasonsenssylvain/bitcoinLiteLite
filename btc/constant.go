package btc

import "github.com/jasoncodingnow/bitcoinLiteLite/crypto"

const PayloadHashSize = 80
const TransactionSignatureSize = 80

const TransactionHeaderSize = crypto.PublicKeyLen + crypto.PublicKeyLen + PayloadHashSize + 4 + 4 + 4

const MaxInt = int(^uint(0) >> 1)
