package utils

import (
	"crypto/md5"
	"encoding/hex"
)

func GetMD5Hash(data string) string {
	hasher := md5.New()
	hasher.Write([]byte(data))
	hashedData := hex.EncodeToString(hasher.Sum(nil))
	return hashedData
}
