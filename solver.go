package main

import (
	"bufio"
	"bytes"
	_ "embed"
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
)

// https://xcloche.hateblo.jp/entry/2022/01/24/212558
// https://slc.is/data/wordles.txt

//go:embed wordles.txt
var wordles_txt []byte

//go:embed oxford-english-dic-words.txt
var oxford_english_dic_words_txt []byte

var score_list []int = make([]int, 243)

var wordle_words []string
var oxford_words []string
var word_list []string

func load_words(src []byte) []string {
	words := make([]string, 0)
	fp := bytes.NewReader(src)
	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		line := scanner.Text()
		words = append(words, line)
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return words
}

func init() {

	for i := 0; i < 243; i++ {
		score_list[i] = (i % 3) + (i/3%3)*10 + (i/9%3)*100 + (i/27%3)*1000 + (i/81%3)*10000
	}

	wordle_words = load_words(wordles_txt)
	oxford_words = load_words(oxford_english_dic_words_txt)
}

func calc_score(answer string, input string) int {
	answer_chars := []rune(answer)
	vs := make([]int, len(input))

	result := 0

	for i, ch := range input {
		if answer_chars[i] == ch {
			vs[i] = 2
			answer_chars[i] = 0
		}
	}

	for i, ch := range input {
		if vs[i] == 0 {
			for j := 0; j < len(answer_chars); j++ {
				if answer_chars[j] == ch {
					vs[i] = 1
					answer_chars[j] = 0
					break
				}
			}
		}
	}

	result = 0
	for i := 0; i < len(vs); i++ {
		result = result*10 + vs[i]
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
	cr := true
	for i, word := range words {
		fmt.Print(word)
		cr = false
		if i%20 == 19 {
			fmt.Print("\n")
			cr = true
		} else {
			fmt.Print(" ")
		}
	}
	if !cr {
		fmt.Print("\n")
	}
}

func filter_by_args(words []string, args []string) []string {
	for i := 0; i+1 < len(args); i += 2 {
		input := args[i]
		output, _ := strconv.Atoi(args[i+1])
		words = filter(words, input, output)
	}
	return words
}

func calc_cmd(args []string) {
	words := filter_by_args(word_list, args)
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
	fmt.Println("Count:", len(words))
	if len(words) <= 20 {
		print_words(words)
	}
}

func solve_cmd(args []string) {
	words := filter_by_args(word_list, args)
	print_words(words)
	fmt.Println("Count:", len(words))
}

func score_cmd(args []string) {
	if len(args) < 2 {
		fmt.Println("Usage:", os.Args[0], "score <answer> <input>")
		os.Exit(1)
	}
	score := calc_score(args[0], args[1])
	fmt.Printf("answer: %s\n", args[0])
	fmt.Printf("input : %s\n", args[1])
	fmt.Printf("score : %05d\n", score)
}

func usage() {
	fmt.Printf("Usage: %s [OPTIONS] COMMAND\n", os.Args[0])
	fmt.Println()
	fmt.Println("  Wordle Solver -> https://www.powerlanguage.co.uk/wordle/")
	fmt.Printf("  dictionary: wordle %d words, oxford %d words\n", len(wordle_words), len(oxford_words))
	fmt.Println()
	fmt.Println("Options:")
	flag.PrintDefaults()
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  score  <ansewer> <input>   : Calculate score")
	fmt.Println("  calc   [<input> <output>]* : Calculate scores of words")
	fmt.Println("  solve  [<input> <output>]* : Print available words set")
	fmt.Println()
	os.Exit(1)
}

func main() {
	var words_flag string
	flag.StringVar(&words_flag, "words", "oxford", "words set (wordle or oxford)")
	flag.Parse()

	if flag.NArg() < 1 {
		usage()
	}

	switch words_flag {
	case "wordle":
		word_list = wordle_words
	case "oxford":
		word_list = oxford_words
	default:
		usage()
	}

	switch flag.Arg(0) {
	case "calc":
		calc_cmd(flag.Args()[1:])
	case "solve":
		solve_cmd(flag.Args()[1:])
	case "score":
		score_cmd(flag.Args()[1:])
	default:
		usage()
	}
}
