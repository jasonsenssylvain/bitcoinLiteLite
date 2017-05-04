package consensus

import (
	"reflect"
)

//CheckProofOfWork 确认POW是否匹配
func CheckProofOfWork(prefixToMatch []byte, hash []byte) bool {
	if len(prefixToMatch) > 0 {
		return reflect.DeepEqual(prefixToMatch, hash[:len(prefixToMatch)])
	}
	return true
}
