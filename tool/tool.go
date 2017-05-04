package tool

import "crypto/sha256"

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

func SliceByteWhenEncount(d []byte, encount byte) []byte {

	for i, bb := range d {

		if bb != encount {
			return d[i:]
		}
	}

	return nil
}
