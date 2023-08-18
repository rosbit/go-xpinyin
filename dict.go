package xpinyin

import (
	"embed"
	"bufio"
	"strings"
	"strconv"
)

const (
	vowelsPattern = "aoeiuv\u00fc"
	xpinyinDictFile = "dict/Mandarin.dat"
)

//go:embed dict
var dictDir embed.FS
var (
	dict map[rune][]string
	toneMark = map[int][]rune{
		0: []rune(vowelsPattern),
		1: []rune("\u0101\u014d\u0113\u012b\u016b\u01d6\u01d6"),
		2: []rune("\u00e1\u00f3\u00e9\u00ed\u00fa\u01d8\u01d8"),
		3: []rune("\u01ce\u01d2\u011b\u01d0\u01d4\u01da\u01da"),
		4: []rune("\u00e0\u00f2\u00e8\u00ec\u00f9\u01dc\u01dc"),
	}
	vowelIndex = func(v rune) int {
		for i, c := range toneMark[0] {
			if c == v {
				return i
			}
		}
		return -1
	}
	retroflex = map[string]struct{}{
		"ZH": struct{}{},
		"CH": struct{}{},
		"SH": struct{}{},
	}
)

func init() {
	f, err := dictDir.Open(xpinyinDictFile)
	if err != nil {
		panic("dict file not found")
	}
	defer f.Close()

	dict = make(map[rune][]string)
	r := bufio.NewReader(f)
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			break
		}
		if len(line) <= 1 {
			continue
		}
		ss := strings.SplitN(line, "\t", 2)
		if len(ss) != 2 {
			continue
		}
		cb, err := strconv.ParseInt(ss[0], 16, 32)
		if err != nil {
			continue
		}
		pinyins := strings.Fields(ss[1])
		if len(pinyins) == 0 {
			continue
		}
		dict[rune(cb)] = pinyins
	}
}
