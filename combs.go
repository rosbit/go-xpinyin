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
func get_comb_indexes(comb_numbers <-chan int, n int) (n_items int, index <-chan int) {
	// calculate the maximal number of possible combinations
	n_max := 1
	num_options_list := []int{}
	for j := range comb_numbers {
		n_max *= j
		n_items += 1
		num_options_list = append(num_options_list, j)
	}
	if n_max < n {
		n = n_max
	}
	if n == 0 {
		panic("0-length list not allowed")
	}

	combIndex := make(chan int)
	go func() {
		i := n_items - 1
		count := 1
		curr := make([]int, n_items)
		for _, idx := range curr {
			combIndex <- idx
		}

		for count < n {
			curr[i] = (curr[i]+1) % num_options_list[i]
			if curr[i] != 0 {
				for _, idx := range curr {
					combIndex <- idx
				}
				count += 1
				i = n_items - 1 // reset to right-most digit
			} else {
				i -= 1 // try previous (left) digit
			}
		}

		close(combIndex)
	}()

	index = combIndex
	return
}

// Given a list of options per place, returns up to n combinations
// e.g.: [['a'], ['1' ,'2'], ['@']] -> [a1@, a2@]
// For instance, ['1' ,'2'] is the group defining the options for the second place
func get_combs(pyOpts <-chan []string, opt *options) (<-chan string) {
	comb_numbers, py_opts := makeCombNumbers(pyOpts)

	combs := make(chan string)
	go func() {
		nItems, combIndex := get_comb_indexes(comb_numbers, opt.maxnCombinations)
		comb, i := make([]string, nItems), 0
		for idx := range combIndex {
			comb[i] = py_opts[i][idx]
			if i=i+1; i >= nItems {
				combs <- strings.Join(comb, opt.splitter)
				i = 0
			}
		}
		close(combs)
	}()
	return combs
}

func makeCombNumbers(pyOpts <-chan []string) (<-chan int, [][]string) {
	py_opts := [][]string{}
	for po := range pyOpts {
		py_opts = append(py_opts, po)
	}
	comb_numbers := make(chan int)

	go func() {
		for _, o := range py_opts {
			comb_numbers <- len(o)
		}
		close(comb_numbers)
	}()
	return comb_numbers, py_opts
}
