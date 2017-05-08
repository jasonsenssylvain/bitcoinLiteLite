# bitcoinLiteLite
a simple crypto currency implementation like bitcoin

***
* doc
  * [block](https://github.com/jasoncodingnow/bitcoinLiteLite/blob/master/block/README.md "")
  * [crypto](https://github.com/jasoncodingnow/bitcoinLiteLite/blob/master/crypto/README.md "")
  * [consensus](https://github.com/jasoncodingnow/bitcoinLiteLite/blob/master/consensus/README.md "")
  * [network](https://github.com/jasoncodingnow/bitcoinLiteLite/blob/master/network/README.md "")

***

* 说明：
  * 该项目是基于bitcoin的设计来写的，目前v0.1.0目标是完成一个最小的可运行的bitcoin网络
  
***

#### 运行

目前暂定首个节点 端口是 8091,在代码 https://github.com/jasoncodingnow/bitcoinLiteLite/blob/master/main.go L 110，请自行修改为当前机器的局域网地址

```GO
git clone https://github.com/jasoncodingnow/bitcoinLiteLite.git
cd github.com/jasoncodingnow/bitcoinLiteLite

go build .
// 启动第一个节点
./bitcoinLiteLite port 8091
// 查看console可以看到第一个节点的 publicKey， 假设是 8091PUBLISKEY

// 启动第二个节点
./bitcoinLiteLite port 8092

// 测试 由第二个节点生成一笔Transaction，并广播
// 目前，5个Transaction会打包，或者比如1个Transaction，会在1分钟内打包
// 在第二个节点的console输入Transaction命令。第一个参数是要传播给谁，第二个参数是消息是什么
8091PUBLISKEY hi

// 等待一分钟打包

```


