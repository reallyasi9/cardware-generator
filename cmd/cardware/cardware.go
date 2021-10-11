package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/reallyasi9/cardware-generator/pkg/cardware"
	"gonum.org/v1/gonum/stat/combin"
)

var flagWordList string
var flagMinWordLength int
var flagNoSymbols bool
var flagQuotes bool
var flagSpace bool
var flagNoCapitals bool

type cardWord struct {
	cards []rune
	word  string
}

type cardWordList []cardWord

func (c cardWordList) Len() int {
	return len(c)
}

func (c cardWordList) Less(i, j int) bool {
	if len(c[i].cards) != len(c[j].cards) {
		return len(c[i].cards) < len(c[j].cards)
	}
	for k := 0; k < len(c[i].cards); k++ {
		if c[i].cards[k] != c[j].cards[k] {
			return c[i].cards[k] < c[j].cards[k]
		}
	}
	return c[i].word < c[j].word
}

func (c cardWordList) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func init() {
	flag.IntVar(&flagMinWordLength, "m", 4, "minimum number of letters in words")
	flag.BoolVar(&flagNoSymbols, "no-symbols", false, "do not create a symbol table")
	flag.BoolVar(&flagQuotes, "quotes", false, "allow quote characters in symbol table")
	flag.BoolVar(&flagSpace, "space", false, "allow space character in symbol table")
	flag.BoolVar(&flagNoCapitals, "no-capitals", false, "do not create a capital letter table")

	flag.Usage = func() {
		name := filepath.Base(os.Args[0])
		fmt.Fprintf(os.Stderr, "Usage: %s [options] wordlist\nOptions are any of the following:\n", name)
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "Options must preceed positional arguments.\n")
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	flag.Parse()
	wordListFile := flag.Arg(0)
	if wordListFile == "" {
		flag.Usage()
		log.Fatal(fmt.Errorf("word list file not specified"))
	}

	file, err := os.Open(wordListFile)
	if err != nil {
		log.Fatal(fmt.Errorf("word list file '%s' : %v", wordListFile, err))
	}
	defer file.Close()

	wordList := cardware.NewWordList(file, flagMinWordLength)
	log.Printf("read %d words", len(wordList))

	nCards := countCardsNeeded(len(wordList))
	log.Printf("needs %d cards", nCards)
	nCards--
	deckSize := len(cardware.FrenchCards)
	nWords := combin.NumPermutations(deckSize, nCards)
	log.Printf("limiting to %d words with %d cards", nWords, nCards)

	rand.Shuffle(len(wordList), func(i, j int) {
		wordList[i], wordList[j] = wordList[j], wordList[i]
	})

	pg := combin.NewPermutationGenerator(deckSize, nCards)
	perm := make([]int, nCards)
	cwl := make(cardWordList, nWords)
	for iWord := 0; pg.Next() && iWord < nWords; iWord++ {
		perm = pg.Permutation(perm)
		cards := make([]rune, nCards)
		for i, c := range perm {
			cards[i] = cardware.FrenchCards[c]
		}
		cwl[iWord] = cardWord{cards: cards, word: wordList[iWord]}
	}

	sort.Sort(cwl)
	for _, cw := range cwl {
		for _, c := range cw.cards {
			cardStr, err := cardware.TranslateFrench(c)
			if err != nil {
				log.Fatal(fmt.Errorf("card '%c' : %v", c, err))
			}
			fmt.Printf("[%s]", cardStr)
		}
		fmt.Printf(" %s\n", cw.word)
	}
}

func countCardsNeeded(nCombinations int) int {
	cards := 0
	combs := 0
	for combs < nCombinations && cards <= len(cardware.FrenchCards) {
		cards++
		combs = combin.NumPermutations(len(cardware.FrenchCards), cards)
	}
	return cards
}