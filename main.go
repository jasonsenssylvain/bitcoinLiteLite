package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"strings"

	"github.com/jasoncodingnow/bitcoinLiteLite/block"
	"github.com/jasoncodingnow/bitcoinLiteLite/blockchainuser"
	"github.com/jasoncodingnow/bitcoinLiteLite/crypto"
	"github.com/jasoncodingnow/bitcoinLiteLite/message"
	"github.com/jasoncodingnow/bitcoinLiteLite/network"
	"github.com/jasoncodingnow/bitcoinLiteLite/tool"
)

var port = flag.String("port", "0", "server port")

func init() {
	flag.Parse()
}

func main() {
	if len(flag.Args()) != 2 {
		fmt.Println(flag.Args())
		fmt.Println("[ERROR] input error")
		return
	}

	port := flag.Args()[1]

	addresses, _ := tool.GetIpAddress()
	Start(addresses[0], port)
	for {
		str := <-ReadStdin()
		if len(str) < 2 {
			fmt.Println("input error")
			fmt.Println(str)
			continue
		}
		block.BlockChainInstant.TransactionChan <- *CreateTransaction(str[0], str[1])
	}
}

func Start(address, port string) {
	fmt.Println("start at " + address + ":" + port)
	//TODO: 需要从有本地保存key的机制
	key := crypto.GenerateKey()
	blockchainuser.SetKey(key)
	fmt.Println("public key is " + key.PublicKey)

	network.NetworkInstant = network.CreateNetwork(address, port)
	go network.Run(network.NetworkInstant, port)
	for _, node := range RootNodes() {
		network.NetworkInstant.ConnectionChan <- node
	}

	block.BlockChainInstant = block.NewBlockChain()
	go block.BlockChainInstant.Run()

	ReceiveMessage()
}

func CreateTransaction(toAddr, msg string) *block.Transaction {
	t := block.NewTransaction([]byte(blockchainuser.GetKey().PublicKey), []byte(toAddr), []byte(msg))
	t.Header.Nonce = t.GenerateNonce(tool.GenerateBytes(block.TransactionPowPrefix, 0))
	t.Signature = t.Sign(blockchainuser.GetKey().PrivateKey)

	return t
}

//接收到外部节点的消息
func ReceiveMessage() {
	go func() {
		for {
			select {
			case msg := <-network.NetworkInstant.IncomingChan:
				HandleIncomingMessage(msg)
			}
		}
	}()
}

//HandleIncomingMessage 处理消息
func HandleIncomingMessage(m message.Message) {
	switch m.Type {
	case message.MessageTypeSendBlock:
		b := block.NewBlock(nil)
		err := b.UnmarshalBinary(m.Data)
		if err != nil {
			fmt.Println("unmarshal binary to block err: " + err.Error())
			break
		}
		block.BlockChainInstant.BlockChan <- *b
	case message.MessageTypeSendTransaction:
		t := block.Transaction{}
		_, err := t.UnmarshalBinary(m.Data)
		if err != nil {
			fmt.Println("unmarshal binary to transaction err: " + err.Error())
			break
		}
		block.BlockChainInstant.TransactionChan <- t
	}
}

//RootNodes 写死的固定服务器，用于接入，暂时先写死自己的一台
func RootNodes() []string {
	nodes := []string{"192.168.11.101:8091"}
	return nodes
}

//ReadStdin 读取命令行
func ReadStdin() chan []string {
	cb := make(chan []string)
	sc := bufio.NewScanner(os.Stdin)

	go func() {
		for {
			if sc.Scan() {
				line := sc.Text()
				fmt.Println("input is " + line)
				lines := strings.Split(line, " ")
				cb <- lines
			}
		}
	}()

	return cb
}
