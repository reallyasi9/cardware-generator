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
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/reallyasi9/cardware-generator/pkg/cardware"

	crand "crypto/rand"
)

var symbols = []rune{'!', '@', '#', '$', '%', '^', '&', '*', '(', ')', '_', '-', '+', '=', '~', '{', '[', '}', ']', '|', '\\', ':', ';', '<', ',', '>', '.', '?', '/'}
var quotes = []rune{'`', '\'', '"'}

type diceBag struct {
	dice []int
}

func (db diceBag) String() string {
	m := make(map[int]int)
	for _, d := range db.dice {
		m[d]++
	}
	k := make([]int, 0)
	for key := range m {
		k = append(k, key)
	}
	sort.Ints(k)
	s := make([]string, len(k))
	for _, die := range k {
		s = append(s, fmt.Sprintf("%dd%d", m[die], die))
	}
	return strings.Join(s, "+")
}

func (db *diceBag) Set(s string) error {
	re := regexp.MustCompile(`(\d+)?[dD](\d+)`)
	for _, val := range strings.Split(s, "+") {
		m := re.FindStringSubmatch(val)
		if m == nil {
			return fmt.Errorf("invalid dice identifier '%s'", val)
		}
		n := 1
		var err error
		if m[1] != "" {
			n, err = strconv.Atoi(m[1])
			if err != nil {
				return err
			}
		}
		die, err := strconv.Atoi(m[2])
		if err != nil {
			return err
		}
		if db.dice == nil {
			db.dice = make([]int, 0)
		}
		for i := 0; i < n; i++ {
			db.dice = append(db.dice, die)
		}
	}
	return nil
}

var flagWordList string
var flagMinWordLength int
var flagNoSymbols bool
var flagQuotes bool
var flagSpace bool
var flagNoCapitals bool
var flagCards int
var flagDiceBag diceBag

func init() {
	//flag.StringVar(&flagWordList, "wordlist", "", "path to file containing list of valid words")
	flag.StringVar(&flagWordList, "w", "", "path to file containing list of valid words")
	//flag.IntVar(&flagMinWordLength, "minlength", 4, "minimum number of letters in words")
	flag.IntVar(&flagMinWordLength, "m", 4, "minimum number of letters in words")
	flag.BoolVar(&flagNoSymbols, "no-symbols", false, "do not create a symbol table")
	flag.BoolVar(&flagQuotes, "quotes", false, "allow quote characters in symbol table")
	flag.BoolVar(&flagSpace, "space", false, "allow space character in symbol table")
	flag.BoolVar(&flagNoCapitals, "no-capitals", false, "do not create a capital letter table")
	flag.IntVar(&flagCards, "c", 0, "draw this many playing cards to augment randomness")
	flag.Var(&flagDiceBag, "d", "define bag of dice (using [N]dF+[N]dF+... notation)")
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

	device := cardware.NewCombined(flagDiceBag.dice)
	log.Printf("using deck: %v", device.Deck)
	log.Printf("using dice: %v", device.DiceBag)

	// figure out number of draws
	kdraws := device.DiceBag.MaxDraws() + flagCards
	nWords := device.CountDistinctOutcomes(kdraws)
	bigLen := big.NewInt(int64(len(wordList)))

	if !nWords.IsInt64() || nWords.Cmp(bigLen) > 0 {
		nWords = bigLen
	}
	nSubset := nWords.Int64()

	log.Printf("drawing a total of %d words", nSubset)

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
	for i, draw := int64(0), device.NextOutcome(kdraws); i < nSubset && draw != nil; i, draw = i+1, device.NextOutcome(kdraws) {
		names := make([]string, len(draw))
		for j, r := range draw {
			names[j], err = device.Translate(r)
			if err != nil {
				panic(err)
			}
		}
		list[i] = strings.Join(names, "+")
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