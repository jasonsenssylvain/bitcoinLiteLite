### btc的几个内容： transaction, block

***

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