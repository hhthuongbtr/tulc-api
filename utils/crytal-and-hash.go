package utils

import (
	"crypto/md5"
	"encoding/hex"
)

func GetMd5FromString(strToSign string) string {
	hasher := md5.New()
	hasher.Write([]byte(strToSign))
	sig := hex.EncodeToString(hasher.Sum(nil))
	return sig
}
