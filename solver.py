import re
import sys

def calc(base, word):
    r = ""
    for i in range(5):
        c = base[i]
        if c == word[i]:
            r += "2"
        elif c in word:
            r += "1"
        else:
            r += "0"
    return r

#qs = [("EARTH", "01100")]

qs = [x.split("=") for x in sys.argv[1:]]

words = [w.rstrip() for w in open("words.txt").readlines()]
count = 0
for word in words:
    if any(calc(q, word) != a for q, a in qs):
        continue
    print(word)
    count += 1
print("count", count)

