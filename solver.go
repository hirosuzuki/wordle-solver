package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
)

var score_list []int

// https://xcloche.hateblo.jp/entry/2022/01/24/212558

// https://slc.is/data/wordles.txt
var word_list []string

func init() {
	score_list = make([]int, 243)
	for i := 0; i < 243; i++ {
		score_list[i] = (i % 3) + (i/3%3)*10 + (i/9%3)*100 + (i/27%3)*1000 + (i/81%3)*10000
	}

	word_list = make([]string, 0)
	fp, err := os.Open("wordles.txt")
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		line := scanner.Text()
		word_list = append(word_list, line)
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}

	fmt.Println("All Words:", len(word_list))
}

func calc_score(answer string, input string) int {
	result := 0
	for i := 0; i < 5; i++ {
		result *= 10
		ch := input[i]
		if answer[i] == ch {
			result += 2
		} else {
			for j := 0; j < 5; j++ {
				if answer[j] == ch {
					result += 1
					break
				}
			}
		}
	}
	return result
}

func filter(words []string, input string, output int) []string {
	result := make([]string, 0)
	for _, word := range words {
		if calc_score(word, input) == output {
			result = append(result, word)
		}
	}
	return result
}

func check() {
	answer := "earth"
	for _, input := range word_list {
		score := calc_score(answer, input)
		fmt.Println(answer, input, score)
	}
}

func calc_max_results_slow(word_list []string, input string) int {
	result := 0
	for _, score := range score_list {
		r := len(filter(word_list, input, score))
		if r > result {
			result = r
		}
	}
	return result
}

func calc_max_results(word_list []string, input string) int {
	score_count := make(map[int]int)
	for i := range score_list {
		score_count[i] = 0
	}
	result := 0
	for _, word := range word_list {
		score := calc_score(word, input)
		score_count[score] += 1
		if score_count[score] > result {
			result = score_count[score]
		}
	}
	return result
}

func calc(words []string, word_list []string) map[string]int {
	word_opts := make(map[string]int)
	for i, input := range word_list {
		fmt.Printf("calc %d/%d\r", i, len(word_list))
		r := calc_max_results(words, input)
		word_opts[input] = r
	}
	return word_opts
}

func print_words(words []string) {
	for i := 0; i < len(words); i += 1 {
		fmt.Print(words[i])
		if i%10 == 9 {
			fmt.Print("\n")
		} else {
			fmt.Print(" ")
		}
	}
	fmt.Println("\nCount:", len(words))
}

func filter_by_args(words []string, args []string) []string {
	for i := 0; i+1 < len(args); i += 2 {
		input := args[i]
		output, _ := strconv.Atoi(args[i+1])
		words = filter(words, input, output)
	}
	return words
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage:", os.Args[0], "[score|solve|calc]")
		os.Exit(1)
	}
	if os.Args[1] == "calc" {
		words := filter_by_args(word_list, os.Args[2:])
		word_opts := calc(words, word_list)
		word_map := make(map[string]bool)
		for _, word := range words {
			word_map[word] = true
		}
		all_words := make([]string, len(word_list))
		copy(all_words, word_list)
		sort.Slice(all_words, func(i, j int) bool {
			d := word_opts[all_words[i]] - word_opts[all_words[j]]
			if d != 0 {
				return d > 0
			}
			if word_map[all_words[j]] {
				return true
			}
			return false
		})
		for _, word := range all_words {
			fmt.Printf("%s %d", word, word_opts[word])
			if word_map[word] {
				fmt.Print("*")
			}
			fmt.Print("\n")
		}

	}
	if os.Args[1] == "solve" {
		words := filter_by_args(word_list, os.Args[2:])
		print_words(words)
	}
	if os.Args[1] == "score" {
		if len(os.Args) < 4 {
			fmt.Println("Usage:", os.Args[0], "score <answer> <input>")
			os.Exit(1)
		}
		score := calc_score(os.Args[2], os.Args[3])
		fmt.Printf("answer: %s\n", os.Args[2])
		fmt.Printf("input : %s\n", os.Args[3])
		fmt.Printf("score : %05d\n", score)
	}
}
