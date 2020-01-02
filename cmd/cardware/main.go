package main

import (
	"bufio"
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"math/big"
	"math/rand"
	"os"
	"sort"
	"strings"

	"github.com/reallyasi9/cardware-generator/pkg/cardware"

	crand "crypto/rand"
)

var symbols = []rune{'!', '@', '#', '$', '%', '^', '&', '*', '(', ')', '_', '-', '+', '=', '~', '{', '[', '}', ']', '|', '\\', ':', ';', '<', ',', '>', '.', '?', '/'}
var quotes = []rune{'`', '\'', '"'}

var flagWordList string
var flagMinWordLength int
var flagNoSymbols bool
var flagQuotes bool
var flagSpace bool
var flagNoCapitals bool

func init() {
	flag.StringVar(&flagWordList, "wordlist", "", "path to file containing list of valid words")
	flag.StringVar(&flagWordList, "w", "", "path to file containing list of valid words")
	flag.IntVar(&flagMinWordLength, "minlength", 4, "minimum number of letters in words")
	flag.IntVar(&flagMinWordLength, "m", 4, "minimum number of letters in words")
	flag.BoolVar(&flagNoSymbols, "no-symbols", false, "do not create a symbol table")
	flag.BoolVar(&flagQuotes, "quotes", false, "allow quote characters in symbol table")
	flag.BoolVar(&flagSpace, "space", false, "allow space character in symbol table")
	flag.BoolVar(&flagNoCapitals, "no-capitals", false, "do not create a capital letter table")
}

func main() {
	flag.Parse()

	file, err := os.Open(flagWordList)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	wordList := make([]string, 0)
	for scanner.Scan() {
		t := scanner.Text()
		if len(t) >= flagMinWordLength {
			wordList = append(wordList, scanner.Text())
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	log.Printf("read %d words", len(wordList))

	deck := cardware.NewStandardFrenchDeck()

	// figure out number of draws needed
	k := 0
	var nWords *big.Int
	bigLen := big.NewInt(int64(len(wordList)))
	for ; k < deck.MaxDraws(); k++ {
		np := deck.NumPermutations(k)
		if np.Cmp(bigLen) > 0 {
			break
		}
		nWords = np
	}
	k--

	if !nWords.IsInt64() {
		log.Fatalf("too many words: %d", nWords)
	}
	nSubset := nWords.Int64()

	log.Printf("using %d draws (for a total of %d words)", k, nSubset)

	// shuffle and select words from the wordlist
	var src cryptoSource
	rng := rand.New(src)
	rng.Shuffle(len(wordList), func(i, j int) {
		wordList[i], wordList[j] = wordList[j], wordList[i]
	})
	subset := wordList[:nSubset]
	sort.Strings(subset)

	// generate permutations
	list := make([]string, nSubset)
	perm := deck.NextPermutation(k)
	i := 0
	for ; perm != nil; perm = deck.NextPermutation(k) {
		name := make([]string, k)
		var cards strings.Builder
		for j, card := range perm {
			cards.WriteRune(card)
			name[j], _ = deck.Translate(card)
		}
		list[i] = fmt.Sprintf("%s (%s)", cards.String(), strings.Join(name, "+"))
		i++
	}

	// sort by card order
	sort.Strings(list)

	// print card-word listing
	for i, draw := range list {
		fmt.Printf("%s %s\n", draw, subset[i])
	}

	// generate symbols
	if !flagNoSymbols {
		if flagQuotes {
			symbols = append(symbols, quotes...)
		}
		if flagSpace {
			symbols = append(symbols, ' ')
		}
		rng.Shuffle(len(symbols), func(i, j int) {
			symbols[i], symbols[j] = symbols[j], symbols[i]
		})

		fmt.Print("\n ")
		for _, col := range cardware.FrenchColors {
			fmt.Printf("  %c", col)
		}
		fmt.Println()

		for i, val := range cardware.FrenchValues {
			fmt.Printf("%c", val)
			for j := range cardware.FrenchColors {
				fmt.Printf("  %c", symbols[i*len(cardware.FrenchColors)+j])
			}
			fmt.Println()
		}
	}

	// generate capitals
	if !flagNoCapitals {
		fmt.Println()
		blackCap := rng.Float32() < .5
		if blackCap {
			fmt.Println("CAPITAL: B")
		} else {
			fmt.Println("CAPITAL: R")
		}
	}
}

type cryptoSource struct{}

func (s cryptoSource) Seed(seed int64) {}

func (s cryptoSource) Int63() int64 {
	return int64(s.Uint64() & ^uint64(1<<63))
}

func (s cryptoSource) Uint64() (v uint64) {
	err := binary.Read(crand.Reader, binary.BigEndian, &v)
	if err != nil {
		log.Fatal(err)
	}
	return v
}
