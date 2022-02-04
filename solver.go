package main

import (
	"bufio"
	"bytes"
	_ "embed"
	"flag"
	"fmt"
	"io"
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
var all_word_list []string
var answers_word_list []string

func load_words(reader io.Reader) []string {
	words := make([]string, 0)
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		words = append(words, line[:5])
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return words
}

func load_words_from_byte(src []byte) []string {
	return load_words(bytes.NewReader(src))
}

func load_words_from_file(src string) []string {
	fp, err := os.Open(src)
	if err != nil {
		panic(err)
	}
	defer fp.Close()
	return load_words(fp)
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

func calc(words []string, word_list []string, verbose bool) map[string]int {
	word_opts := make(map[string]int)
	for i, input := range word_list {
		if verbose {
			fmt.Printf("calc %d/%d\r", i, len(word_list))
		}
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
	words := filter_by_args(answers_word_list, args)
	word_opts := calc(words, all_word_list, true)
	word_map := make(map[string]bool)
	for _, word := range words {
		word_map[word] = true
	}
	all_words := make([]string, len(all_word_list))
	copy(all_words, all_word_list)
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
	words := filter_by_args(answers_word_list, args)
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

func choice_word(word_list []string, args []string) string {
	if len(args) == 0 {
		return "raise"
	}
	words := filter_by_args(word_list, args)
	word_map := make(map[string]bool)
	for _, word := range words {
		word_map[word] = true
	}
	options := words
	word_opts := calc(words, options, false)
	all_words := make([]string, len(options))
	copy(all_words, options)
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
	return all_words[len(all_words)-1]
}

func sim_cmd(args []string) {
	words := filter_by_args(all_word_list, args)
	hands_map := make([]int, 10)
	for i, answer := range words {
		args := make([]string, 0)
		for k := 0; k < 10; k++ {
			input := choice_word(words, args)
			score := calc_score(answer, input)
			args = append(args, input, strconv.Itoa(score))
			if score == 22222 {
				break
			}
		}
		fmt.Println(i, answer, args)
		hands_map[len(args)/2] += 1
	}
	t := 0
	c := 0
	for i := 0; i < len(hands_map); i++ {
		fmt.Println(i, hands_map[i])
		t += i * hands_map[i]
		c += hands_map[i]
	}
	fmt.Println(t, c, float64(t)/float64(c))
}

func calc2_cmd(args []string) {
	words := filter_by_args(all_word_list, args)
	for i, input1 := range words {
		for _, score := range score_list {
			result1 := filter(words, input1, score)
			max_result := 0
			max_result_word := "aaaaa"
			for _, input2 := range words {
				r := calc_max_results(result1, input2)
				if r > max_result {
					max_result = r
					max_result_word = input2
				}
				// fmt.Println(input1, score, input2, r)
			}
			fmt.Println(input1, score, max_result, max_result_word)

		}
		fmt.Print(i, "/", len(words), " ", input1, "\n")
		break
	}
	// print(words)
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
	fmt.Println("  sim    [<input> <output>]* : Calculate 2 scores of words")
	fmt.Println()
	os.Exit(1)
}

func main() {

	for i := 0; i < 243; i++ {
		score_list[i] = (i % 3) + (i/3%3)*10 + (i/9%3)*100 + (i/27%3)*1000 + (i/81%3)*10000
	}

	wordle_words = load_words_from_byte(wordles_txt)
	oxford_words = load_words_from_byte(oxford_english_dic_words_txt)

	var dic_flag string
	var answers_flag string
	flag.StringVar(&dic_flag, "dic", "oxford", "all words set (wordle or oxford or filename)")
	flag.StringVar(&answers_flag, "answers", "oxford", "answers words set (wordle or oxford or filename)")
	flag.Parse()

	if flag.NArg() < 1 {
		usage()
	}

	switch dic_flag {
	case "wordle":
		all_word_list = wordle_words
	case "oxford":
		all_word_list = oxford_words
	default:
		all_word_list = load_words_from_file(dic_flag)
	}

	switch answers_flag {
	case "wordle":
		answers_word_list = wordle_words
	case "oxford":
		answers_word_list = oxford_words
	default:
		answers_word_list = load_words_from_file(answers_flag)
	}

	switch flag.Arg(0) {
	case "calc":
		calc_cmd(flag.Args()[1:])
	case "solve":
		solve_cmd(flag.Args()[1:])
	case "score":
		score_cmd(flag.Args()[1:])
	case "sim":
		sim_cmd(flag.Args()[1:])
	case "calc2":
		calc2_cmd(flag.Args()[1:])
	default:
		usage()
	}
}
