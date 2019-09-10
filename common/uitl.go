package common

import "github.com/axgle/mahonia"

//GBK2UTF8 gbk conver utf-8
func GBK2UTF8(gbk string) string {
	return mahonia.NewDecoder("gbk").ConvertString(gbk)
}
