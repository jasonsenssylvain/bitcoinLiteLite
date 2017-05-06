### btc的几个内容： transaction, block

***

#### 流程

* BlockChain，有个属性 Block，该属性放的是当前持有的Block，还未写入到Chain里。当符合一定条件，目前设定是5个Block或者十分钟到了数量还未到5个，就打包。
* 打包的时候有个注意事项，因为打包的BLock可能是自己产生的，也可能是别人产生的，所以很可能打包进入的Block里的Transaction与自己当前的Block里的Transaction不相同，所以需要取出不同放到下一个Block

##### Transaction

* Header: 
  * From        从哪个地址发出
  * To          发到哪个地址
  * PayloadHash 对交易内容进行SHA256
  * PayloadLen  交易内容的长度，用于最终解码
  * Timestamp   交易产生时间
  * Nonce       可以理解为不断递增的数字，用于跟交易内容合并，产生HASH，然后匹配是否符合POW

* Transaction:
  * Header      交易Header
  * Signature   交易发起者使用私钥进行签名
  * Payload     交易内容

#### Block

* BlockHeader
  * Origin      打包者的地址/公钥
  * PrevBlock   上一个Block的地址
  * MerkleRoot  包含的交易形成的Merkle
  * Timestamp   打包时间
  * Nonce       可以理解为不断递增的数字，用于跟交易内容合并，产生HASH，然后匹配是否符合POW

* Block
  * Header        包的具体内容
  * Signature     打包者用自己私钥对当前的block的Hash进行签名
  * Transactions  该区块包含的所有transaction 

#### BlockChain

* BlockChain
  * Block 当前所持有的Block，还未写入到BlockSlice里
  * BlockSlice 当前的Blockchain
  * BlockChan 接收外界产生的以及自己产生的Block的channel
  * TransactionChan 接收外界或者自己产生的Transaction