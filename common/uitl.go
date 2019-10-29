package common

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"

	"github.com/axgle/mahonia"
)

//GBK2UTF8 gbk conver utf-8
func GBK2UTF8(gbk string) string {
	return mahonia.NewDecoder("gbk").ConvertString(gbk)
}

//MD5 str to md5
func MD5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

//EncodeStdBase64 encode str to base64
func EncodeStdBase64(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

//DecodeStdBase64 decode base64 string to string
func DecodeStdBase64(str string) string {
	source, _ := base64.StdEncoding.DecodeString(str)
	return string(source)
}

//EncodeURLBase64 encode str to base64
func EncodeURLBase64(str string) string {
	return base64.URLEncoding.EncodeToString([]byte(str))
}

//DecodeURLBase64 decode base64 string to string
func DecodeURLBase64(str string) string {
	source, _ := base64.URLEncoding.DecodeString(str)
	return string(source)
}
