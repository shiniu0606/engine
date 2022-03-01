package base

import (
	"crypto/md5"
	//"encoding/gob"
	"encoding/hex"
)

var (
	md5_salt = "0fdfa5e5a88gefae640k5d88e7c84708"
)

func MD5Str(s string) string {
	return MD5Bytes([]byte(s))
}

func MD5Bytes(s []byte) string {
	md5Ctx := md5.New()
	md5Ctx.Write(s)
	cipherStr := md5Ctx.Sum(nil)
	return hex.EncodeToString(cipherStr)
}

func MD5WithSalt(s string) string {
	return MD5Bytes([]byte(s + md5_salt))
}
