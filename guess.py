import sys

VOWEL_CHARS = "aiueo"

past_words = [line[:5] for line in open("wordle-history.txt")]

for line in sys.stdin:
    line = line.strip()
    if line.startswith("#"):
        continue
    if len(line) < 5:
        continue
    word = line[:5]
    score = 10
    vowel_count = len([c for c in word if c in VOWEL_CHARS])
    if vowel_count in [0, 3, 4, 5]:
        score = 0
    if word[0] in VOWEL_CHARS:
        score = 0
    if word[4] in VOWEL_CHARS:
        score = 0
    if word in past_words:
        score = 0
    # print(word, score, vowel_count)
    if score:
        print(word)


