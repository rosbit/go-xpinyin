package xpinyin

import (
	"strings"
)

// Given a list with the number of possible options per place, returns a list of numbers representing combinations.
// The combinations are created via additions to a multi-radix number, from left to right
// (i.e. from smaller to larger numbers).
//
// e.g. [2, 2, 1] -> [[0, 0, 0], [0, 1, 0], [1, 0, 0], [1, 1, 0]]
// i.e. we have 2 options (0, 1) for the first and second places and one option (0) for the third.
func get_comb_indexes(num_options_list []int, n int) (<-chan []int) {
	// calculate the maximal number of possible combinations
	n_max := 1
	for _, j := range num_options_list {
		n_max *= j
	}
	if n_max < n {
		n = n_max
	}
	if n == 0 {
		panic("0-length list not allowed")
	}

	n_items := len(num_options_list)
	combs := make(chan []int)
	go func() {
		i := n_items - 1
		count := 1
		curr := make([]int, n_items)
		currCopy := make([]int, n_items)
		copy(currCopy, curr)
		combs <- currCopy

		for count < n {
			curr[i] = (curr[i]+1) % num_options_list[i]
			if curr[i] != 0 {
				currCopy = make([]int, n_items)
				copy(currCopy, curr)
				combs <- currCopy
				count += 1
				i = n_items - 1 // reset to right-most digit
			} else {
				i -= 1 // try previous (left) digit
			}
		}

		close(combs)
	}()

	return combs
}

// Given a list of options per place, returns up to n combinations
// e.g.: [['a'], ['1' ,'2'], ['@']] -> [a1@, a2@]
// For instance, ['1' ,'2'] is the group defining the options for the second place
func get_combs(py_opts [][]string, opt *options) (<-chan string) {
	combs := make(chan string)
	comb_numbers := make([]int, len(py_opts))
	for i, o := range py_opts {
		comb_numbers[i] = len(o)
	}
	combs_indexes := get_comb_indexes(comb_numbers, opt.maxnCombinations)

	go func() {
		for c := range combs_indexes {
			comb := []string{}
			for i:=0; i<len(c); i++ {
				comb = append(comb, py_opts[i][c[i]])
			}
			combs <- strings.Join(comb, opt.splitter)
		}
		close(combs)
	}()
	return combs
}

