package cardware

import (
	"fmt"
	"math/big"

	"gonum.org/v1/gonum/stat/combin"
)

// Card represents a single card in a deck of playing cards.
type Card rune

// Deck represents a deck of playing cards.
type Deck struct {
	cards []Card
	draws int
	tr    func(rune) (string, error)
	pg    *combin.PermutationGenerator
	cg    *combin.CombinationGenerator
}

const AceOfSpades = '🂡'

var FrenchColors = []rune{'B', 'R'}
var FrenchSuits = []rune{'♠', '♡', '♢', '♣'}
var FrenchValues = []rune{'A', '2', '3', '4', '5', '6', '7', '8', '9', 'T', 'J', 'Q', 'K'}
var FrenchCards = make([]rune, 52)

func init() {
	for i := 0; i < 52; i++ {
		FrenchCards[i] = rune(AceOfSpades + (i/13)*16 + (i % 13))
	}
}

// NewStandardFrenchDeck builds a 52-card deck with four French suits
// (clubs, hearts, diamonds, spades) and thirteen French values (ace
// through king)
func NewStandardFrenchDeck() *Deck {
	cards := make([]Card, 52)
	for i, c := range FrenchCards {
		cards[i] = Card(c)
	}
	return &Deck{cards: cards, draws: -1, tr: TranslateFrench}
}

// MaxDraws implements RandomObject interface.
func (d *Deck) MaxDraws() int {
	return len(d.cards)
}

// NumPermutations implements RandomObject interface.
func (d *Deck) NumPermutations(k int) *big.Int {
	md := d.MaxDraws()
	if k > md {
		panic("k > MaxDraws")
	}
	if k < 0 {
		panic("k < 0")
	}
	if k == 0 {
		return big.NewInt(1)
	}
	n := big.NewInt(0)
	return n.MulRange(int64(md-k+1), int64(md))
}

// NumCombinations implements RandomObject interface.
func (d *Deck) NumCombinations(k int) *big.Int {
	md := d.MaxDraws()
	if k > md {
		panic("k > MaxDraws")
	}
	if k < 0 {
		panic("k < 0")
	}
	if k == 0 {
		return big.NewInt(1)
	}
	// (n,k) = (n,n-k)
	if k > md/2 {
		k = md - k
	}
	n := big.NewInt(1)
	n.MulRange(int64(md-k+1), int64(md))
	denom := big.NewInt(1)
	denom.MulRange(int64(1), int64(k))
	return n.Quo(n, denom)
}

// NextPermutation implements RandomObject interface.
func (d *Deck) NextPermutation(k int) []rune {
	md := d.MaxDraws()
	if k > md {
		panic("k > MaxDraws")
	}
	if k < 0 {
		panic("k < 0")
	}
	if k == 0 {
		return []rune{}
	}
	if k != d.draws || d.pg == nil {
		d.draws = k
		d.pg = combin.NewPermutationGenerator(md, k)
	}
	if !d.pg.Next() {
		d.pg = nil
		return nil
	}
	p := make([]int, k)
	d.pg.Permutation(p)
	out := make([]rune, k)
	for i, x := range p {
		out[i] = rune(d.cards[x])
	}
	return out
}

// NextCombination implements RandomObject interface.
func (d *Deck) NextCombination(k int) []rune {
	md := d.MaxDraws()
	if k > md {
		panic("k > MaxDraws")
	}
	if k < 0 {
		panic("k < 0")
	}
	if k == 0 {
		return []rune{}
	}
	if k != d.draws || d.cg == nil {
		d.draws = k
		d.cg = combin.NewCombinationGenerator(md, k)
	}
	if !d.cg.Next() {
		d.cg = nil
		return nil
	}
	c := make([]int, k)
	d.cg.Combination(c)
	out := make([]rune, k)
	for i, x := range c {
		out[i] = rune(d.cards[x])
	}
	return out
}

// TranslateFrench translates a playing card rune into a text name.
func TranslateFrench(r rune) (string, error) {
	suit := int(r-AceOfSpades) / 16
	if suit < 0 || suit >= len(FrenchSuits) {
		return "", fmt.Errorf("suit of card '%c' is unknown", r)
	}
	value := int(r-AceOfSpades) % 16
	if value < 0 || value >= len(FrenchValues) {
		return "", fmt.Errorf("value of card '%c' is unknown", r)
	}
	return string(FrenchValues[value]) + string(FrenchSuits[suit]), nil
}

// Translate implements RandomObject interface.
func (d *Deck) Translate(r rune) (string, error) {
	return d.tr(r)
}
