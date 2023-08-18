// Translate Chinese hanzi to pinyin (拼音) by Python, 汉字转拼音

package xpinyin

import (
	"strings"
	"bytes"
	"fmt"
	"regexp"
)

var (
	reVowels *regexp.Regexp
)

func init() {
	var err error
	reVowels, err = regexp.Compile(fmt.Sprintf("[%s]+", vowelsPattern))
	if err != nil {
		panic(err)
	}
}

func exists(vv string, char_py_options []string) bool {
	for _, o := range char_py_options {
		if strings.HasPrefix(o, vv) {
			return true
		}
	}
	return false
}

// Get All pinyin combinations given all possible readings of each character.
// The number of combinations is limited par default to 10 to avoid exponential explosion on long texts.
func GetPinyins(s string, opts ...Option) []string {
	res := []string{}

	for r := range getPinyins(s, opts...) {
		res = append(res, r)
	}
	return res
}

func getPinyins(s string, opts ...Option) (<-chan string) {
	op := getOptions(opts...)
	all_pinyin_options := [][]string{} // a list of lists that we'll fill with all pinyin options for each character
	is_in_list := true                 // in the list (otherwise, probably not a Chinese character)
	var lastOpts []string
	var lastSeq *string
	for _, ch := range s {
		pinyins, ok := dict[ch]
		if !ok {
			if is_in_list {
				is_in_list = false // within a sequence of non Chinese characters
				lastOpts = []string{toStr(ch)} // add as is
				lastSeq = &lastOpts[len(lastOpts)-1]
			} else {
				*lastSeq = *lastSeq + toStr(ch) // add to previous sequence of non Chinese chars
			}
		} else {
			if lastOpts != nil {
				all_pinyin_options = append(all_pinyin_options, lastOpts)
				lastOpts, lastSeq = nil, nil
			}

			var char_py_options []string
			if !op.toneMarks && !op.toneNumbers {
				// in this case we may have duplicates if the variations differ just by the tones
				for _, v := range pinyins {
					vv := v[:len(v)-1]
					if !exists(vv, char_py_options) {
						char_py_options = append(char_py_options, vv)
					}
				}
			} else {
				char_py_options = pinyins
			}

			var last int
			if op.maxnCombinations == 1 {
				last = 1
			} else {
				last = len(char_py_options)
			}

			var char_options []string
			if op.toneMarks {
				char_options = make([]string, last)
				for i, o := range char_py_options[:last] {
					char_options[i] = decode_pinyin(o)
				}
			} else {
				// 'numbers' or None
				char_options = char_py_options[:last]
			}

			c := len(char_options)
			pys := make([]string, c)
			for i, o := range char_options {
				pys[i] = convert_pinyin(o, op)
			}
			all_pinyin_options = append(all_pinyin_options, pys)
			is_in_list = true
		}
	}
	if lastOpts != nil {
		all_pinyin_options = append(all_pinyin_options, lastOpts)
		lastOpts, lastSeq = nil, nil
	}
	return get_combs(all_pinyin_options, op)
}

func decode_pinyin(s string) (res string) {
	s = strings.ToLower(s)

	r := &bytes.Buffer{}
	t := &bytes.Buffer{}
	for _, c := range s {
		if 'a' <= c && c <= 'z' {
			t.WriteRune(c)
		} else if c == ':' {
			l := t.Len()
			if l >= 1 {
				t.Truncate(l-1)
				t.WriteRune('\u00fc')
			}
		} else {
			if '0' <= c && c <= '5' {
				if tone := int((c - '0') % 5); tone != 0 {
					ts := t.String()
					loc := reVowels.FindStringIndex(ts)
					if len(loc) == 0 {
						// pass when no vowels find yet
						t.WriteRune(c)
					} else if runes := []rune(ts[loc[0]:loc[1]]); len(runes) == 1 {
						// if just find one vowels, put the mark on it
						t.Reset()
						t.WriteString(ts[:loc[0]])
						t.WriteRune(toneMark[tone][vowelIndex(runes[0])])
						t.WriteString(ts[loc[1]:])
					} else {
						// mark on vowels which search with "a, o, e" one by one
						// when "i" and "u" stand together, make the vowels behind
						for i, vowel := range []string{"a", "o", "e", "ui", "iu"} {
							if strings.Index(ts, vowel) >= 0 {
								ts = strings.Replace(ts, toStr(rune(vowel[len(vowel)-1])), toStr(toneMark[tone][i]), 1)
								t.Reset()
								t.WriteString(ts)
								break
							}
						}
					}
				}
			}
			r.WriteString(t.String())
			t.Reset()
		}
	}
	r.WriteString(t.String())
	res = r.String()
	return
}

func convert_pinyin(o string, opt *options) string {
	if opt.toCapitalize {
		return strings.Title(o)
	}
	if opt.toUpper {
		return strings.ToUpper(o)
	}
	return strings.ToLower(o)
}

func GetPinyin(s string, opts ...Option) (res string) {
	opts = append(opts, MaxNCombinations(1))
	pinyins := getPinyins(s, opts...)
	res = <-pinyins
	go func() {
		for range pinyins {
			// discard it
		}
	}()
	return
}
var toStr = func(ch rune) string {return fmt.Sprintf("%c", ch)}
var getInitial = func(ch rune, withRetroflex bool) (res string) {
	pinyin, ok := dict[ch]
	if !ok || len(pinyin) == 0 {
		return toStr(ch)
	}
	firstPinyin := pinyin[0]

	if !withRetroflex {
		return firstPinyin[:1]
	}
	if len(firstPinyin) >= 2 {
		if _, ok := retroflex[firstPinyin[:2]]; ok {
			return firstPinyin[:2]
		}
	}
	return firstPinyin[:1]
}

func GetInitial(ch rune, opts ...Option) (res string) {
	op := getOptions(opts...)
	return getInitial(ch, op.retroflex)
}

func GetInitials(s string, opts ...Option) (res string) {
	op := getOptions(opts...)
	initials := getInitials(s, op.retroflex)
	sb := &strings.Builder{}
	sb.WriteString(<-initials)
	for initial := range initials {
		sb.WriteString(op.splitter)
		sb.WriteString(initial)
	}
	return sb.String()
}

func getInitials(s string, withRetroflex bool) (<-chan string) {
	initials := make(chan string)
	go func() {
		for _, ch := range s {
			initials <- getInitial(ch, withRetroflex)
		}
		close(initials)
	}()
	return initials
}
