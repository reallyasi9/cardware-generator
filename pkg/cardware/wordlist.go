package cardware

import (
	"bufio"
	"io"
	"strings"
)

// WordList repreents a parsed word list file.
type WordList []string

func NewWordList(r io.Reader, minWordLength int) []string {
	scanner := bufio.NewScanner(r)
	wordList := make([]string, 0)
	for scanner.Scan() {
		t := strings.TrimSpace(scanner.Text())
		if len(t) > 0 && len(t) >= minWordLength {
			wordList = append(wordList, t)
		}
	}
	return wordList
}
