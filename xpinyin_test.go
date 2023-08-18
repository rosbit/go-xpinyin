package xpinyin

import (
	"testing"
	"fmt"
)

func TestPinyin(t *testing.T) {
	testGetPinyin("上海")
	testGetInitial('上')
	testGetInitials("上海")
	testGetPinyins("模型")
	testGetPinyins("模样")
}

func testGetPinyin(s string) {
	fmt.Printf("GetPinyin(\"%s\")\n => %s\n", s, GetPinyin(s))
	fmt.Printf("GetPinyin(\"%s\", WithToneMarks())\n => %s\n", s, GetPinyin(s, WithToneMarks()))
	fmt.Printf("GetPinyin(\"%s\", WithToneNumbers())\n => %s\n", s, GetPinyin(s, WithToneNumbers()))
	fmt.Printf("GetPinyin(\"%s\", WithSplitter(\"\"))\n => %s\n", s, GetPinyin(s, WithSplitter("")))
	fmt.Printf("GetPinyin(\"%s\", WithSplitter(\" \"))\n => %s\n", s, GetPinyin(s, WithSplitter(" ")))
}

func testGetInitial(ch rune) {
	fmt.Printf("GetInitial('%c')\n => %s\n", ch, GetInitial(ch))
}

func testGetInitials(s string) {
	fmt.Printf("GetInitials(\"%s\")\n => %s\n", s, GetInitials(s))
	fmt.Printf("GetInitials(\"%s\", WithSplitter(\"\"))\n => %s\n", s, GetInitials(s, WithSplitter("")))
	fmt.Printf("GetInitials(\"%s\", WithSplitter(\" \"))\n => %s\n", s, GetInitials(s, WithSplitter(" ")))
	fmt.Printf("GetInitials(\"%s\", WithRetroflex())\n => %s\n", s, GetInitials(s, WithRetroflex()))
}

func testGetPinyins(s string) {
	fmt.Printf("GetPinyins(\"%s\", WithToneMarks())\n => %v\n", s, GetPinyins(s, WithToneMarks()))
}
