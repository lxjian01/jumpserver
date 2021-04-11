package utils

import (
	"fmt"
	"testing"
)

func TestConvert(t *testing.T) {
	s := "我的测试"
	gbk, err := Utf8ToGbk([]byte(s))
	if err != nil {
		fmt.Errorf("Test convert utf-8 to gbk error %v \n",err)
	} else {
		fmt.Println("以GBK的编码方式打开:", string(gbk))
	}

	utf8, err := GbkToUtf8(gbk)
	if err != nil {
		fmt.Errorf("Test convert gbk to utf-8 error %v \n",err)
	} else {
		fmt.Println("以UTF8的编码方式打开:", string(utf8))
	}
}
