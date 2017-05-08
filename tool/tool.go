package tool

import (
	"crypto/rand"
	"crypto/sha256"
	"net"
	"os"
	"time"
)

//FillBytesToFront 把数据截取/填充到指定长度
func FillBytesToFront(data []byte, totalLen int) []byte {
	if len(data) < totalLen {
		delta := totalLen - len(data)
		appendByte := []byte{}
		for delta != 0 {
			appendByte = append(appendByte, 0)
			delta--
		}
		return append(appendByte, data...)
	}
	return data[:totalLen]
}

//SHA256 对数据进行SHA256 Hash
func SHA256(data []byte) []byte {
	hash := sha256.New()
	hash.Write(data)
	return hash.Sum(nil)
}

//SliceByteWhenEncount 如果遇到了
func SliceByteWhenEncount(d []byte, encount byte) []byte {
	for i, bb := range d {
		if bb != encount {
			return d[i:]
		}
	}
	return nil
}

//RandomString 产生随机字符串
func RandomString(n int) string {
	alphanum := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, n)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}
	return string(bytes)
}

//GenerateBytes 根据长度和每个byte的数据，产生bytes
func GenerateBytes(length int, b byte) []byte {
	bytes := []byte{}
	for length != 0 {
		bytes = append(bytes, b)
		length--
	}
	return bytes
}

//Timeout 定时任务，时间到触发channel
func Timeout(t time.Duration) chan bool {
	i := make(chan bool)
	go func() {
		time.Sleep(t)
		i <- true
	}()
	return i
}

//GetIpAddress 获取当前ip
func GetIpAddress() ([]string, error) {
	name, err := os.Hostname()
	if err != nil {
		return nil, err
	}
	addrs, err := net.LookupHost(name)
	if err != nil {
		return nil, err
	}
	return addrs, nil
}
