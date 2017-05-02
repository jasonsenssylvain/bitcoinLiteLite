package tool

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
