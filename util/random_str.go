package util

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func GenerateRandomStr(size int) string {
	bts := make([]byte, size)
	if _, err := rand.Read(bts); err != nil {
		fmt.Println(err)
	}
	return base64.URLEncoding.EncodeToString(bts)
}
