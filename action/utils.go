package action

import (
	"bytes"
	"crypto/rsa"
	"strings"

	"gitlab.ghn.vn/internal-tools/message/cip"
)

func DecryptMsg(msg []byte, key *rsa.PrivateKey) string {

	if strings.Contains(string(msg), "http") {
		return string(msg)
	}

	bb := bytes.Buffer{}
	bb.Write(msg)
	byteArr := cip.DecryptData(key, bb.Bytes())
	return string(byteArr)
}

func GetReceiverId(users []string, sender string) string {

	for _, u := range users {
		if u != sender {
			return u
		}
	}

	return "ERR_RECEIVER"
}

func RemoveUser(arr []string, e string) []string {
	index := -1
	for i, element := range arr {
		if element == e {
			// result = append(result, element)
			index = i
			break
		}
	}
	if index == -1 {
		return nil
	}
	return append(arr[:index], arr[index+1:]...)
}

func IsExist(arr []string, e string) bool {
	for _, element := range arr {
		if element == e {
			return true
		}
	}
	return false
}
