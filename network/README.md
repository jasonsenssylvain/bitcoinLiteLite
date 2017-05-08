* Node
  * Net         连接实例
  * ConnectTime 获取连接的时间

* Network
  * Nodes               当前节点连接的所有外部节点
  * ConnectionChan      新的连接接入，通过该chan触发处理机制
  * Address             当前节点的地址
  * NodeCallback        新的节点连接完成，调用该chan
  * BroadCastChan       当前节点产生的Message丢入该chan，该chan会广播
  * IncomingChan        如果接收到外部消息，丢入到该chan