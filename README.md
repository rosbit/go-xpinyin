# go-xpinyin

Translate Chinese hanzi to pinyin (拼音) by Golang, 汉字转拼音.

Inspired by https://github.com/lxneng/xpinyin

## Usage

```go
package main

import (
    p "github.com/rosbit/go-xpinyin"
    "fmt"
)

func main() {
    testGetPinyin("上海")
}

func testGetPinyins(s string) {
    fmt.Printf("GetPinyin(\"%s\")\n => %s\n", s, p.GetPinyin(s))
    fmt.Printf("GetPinyin(\"%s\", WithToneMarks())\n => %s\n", s, p.GetPinyin(s, p.WithToneMarks()))
    fmt.Printf("GetPinyin(\"%s\", WithToneNumbers())\n => %s\n", s, p.GetPinyin(s, p.WithToneNumbers()))
    fmt.Printf("GetPinyin(\"%s\", WithSplitter(\"\"))\n => %s\n", s, p.GetPinyin(s, p.WithSplitter("")))
    fmt.Printf("GetPinyin(\"%s\", WithSplitter(\" \"))\n => %s\n", s, p.GetPinyin(s, p.WithSplitter(" ")))
}
```

## More

To see the detailed usage, run `go test`, the output is:

```
GetPinyin("上海")
 => shang-hai
GetPinyin("上海", WithToneMarks())
 => shàng-hǎi
GetPinyin("上海", WithToneNumbers())
 => shang4-hai3
GetPinyin("上海", WithSplitter(""))
 => shanghai
GetPinyin("上海", WithSplitter(" "))
 => shang hai
GetInitial('上')
 => S
GetInitials("上海")
 => S-H
GetInitials("上海", WithSplitter(""))
 => SH
GetInitials("上海", WithSplitter(" "))
 => S H
GetInitials("上海", WithRetroflex())
 => SH-H
GetPinyins("模型", WithToneMarks())
 => [mó-xíng mú-xíng]
GetPinyins("模样", WithToneMarks())
 => [mó-yáng mó-yàng mó-xiàng mú-yáng mú-yàng mú-xiàng]
PASS
ok  	github.com/rosbint/go-xpinyin	0.570s
```
