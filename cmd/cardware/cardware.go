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

var flagMinWordLength int
var flagDraws int
var flagDeckType string

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
	flag.IntVar(&flagDraws, "n", 0, "number of card draws (limits the number of shuffled words; defaults to as many as necessary to select all words in wordlist)")
	flag.StringVar(&flagDeckType, "t", "french", "type of deck (can be \"french\" for a standard 4-suited, 13-ranked deck or \"tarot\" for a 4-suited, 14-ranked, 22-trump deck)")

	flag.Usage = func() {
		name := filepath.Base(os.Args[0])
		fmt.Fprintf(os.Stderr, "Usage: %s [options] wordlist\nOptions are any of the following:\n", name)
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "Options must precede positional arguments.\n")
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
	var deck *cardware.Deck
	if flagDeckType == "french" {
		deck = cardware.NewStandardFrenchDeck()
	} else if flagDeckType == "tarot" {
		deck = cardware.NewTarotDeMarseilleDeck()
	} else {
		flag.Usage()
		log.Fatal(fmt.Errorf("deck type \"%s\" not valid", flagDeckType))
	}

	file, err := os.Open(wordListFile)
	if err != nil {
		log.Fatal(fmt.Errorf("word list file '%s' : %v", wordListFile, err))
	}
	defer file.Close()

	wordList := cardware.NewWordList(file, flagMinWordLength)
	log.Printf("read %d words", len(wordList))

	nCards := countCardsNeeded(len(wordList), deck)
	log.Printf("needs %d cards", nCards)
	if nCards > flagDraws && flagDraws > 0 {
		log.Printf("limiting to %d cards due to user options", flagDraws)
		nCards = flagDraws
	}
	deckSize := deck.MaxDraws()
	nWords := combin.NumPermutations(deckSize, nCards)
	if nWords > len(wordList) {
		log.Printf("WARNING: due to wordlist size, only %d of %d permutations will be used", len(wordList), nWords)
		nWords = len(wordList)
	}
	log.Printf("limiting to %d words with %d cards", nWords, nCards)

	rand.Shuffle(len(wordList), func(i, j int) {
		wordList[i], wordList[j] = wordList[j], wordList[i]
	})

	cwl := make(cardWordList, nWords)
	for iWord := 0; iWord < nWords; iWord++ {
		cards := deck.NextOutcome(nCards)
		cwl[iWord] = cardWord{cards: cards, word: wordList[iWord]}
	}

	sort.Sort(cwl)
	for _, cw := range cwl {
		for _, c := range cw.cards {
			cardStr, err := deck.Translate(c)
			if err != nil {
				log.Fatal(fmt.Errorf("card '%c' : %v", c, err))
			}
			fmt.Printf("[%s]", cardStr)
		}
		fmt.Printf(" %s\n", cw.word)
	}
}

func countCardsNeeded(nCombinations int, deck *cardware.Deck) int {
	cards := 0
	combs := 0
	for combs < nCombinations && cards <= deck.MaxDraws() {
		cards++
		combs = combin.NumPermutations(deck.MaxDraws(), cards)
	}
	return cards
}
