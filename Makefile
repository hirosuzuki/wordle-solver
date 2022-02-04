wordle-solver: solver.go
	go build -o wordle-solver

first-calc: wordle-solver guess-words.txt oxford-english-dic-words.txt
	./wordle-solver -dic oxford-english-dic-words.txt -answers guess-words.txt calc

guess-words.txt: oxford-english-dic-words.txt wordle-history.txt
	python3 guess.py < oxford-english-dic-words.txt > guess-words.txt

clean:
	rm wordle-solver