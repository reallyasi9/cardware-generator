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
	RandomObject
	cards []Card
	draws int
	tr    func(rune) (string, error)
	pg    *combin.PermutationGenerator
}

// AceOfSpades is the lowest valued card in the deck.
const AceOfSpades = '🂡'

// KingOfClubs is the highest valued non-trump card in the deck.
const KingOfClubs = '🃞'

// TheFool is the lowest valued trump card in the deck.
const TheFool = '🃠'

// TheWorld is the highest valued trump card in the deck.
const TheWorld = '🃵'

// FrenchColors are the colors present in a standard French deck of cards (black and red).
var FrenchColors = []rune{'B', 'R'}

// FrenchSuits are the four suits present in a standard French deck of cards.
var FrenchSuits = []rune{'♠', '♡', '♢', '♣'}

// FrenchValues are the thirteen values of cards for each suit in a standard French deck of cards.
var FrenchValues = []rune{'A', '2', '3', '4', '5', '6', '7', '8', '9', 'T', 'J', 'Q', 'K'}

// FrenchCards is a deck of 52 standard French cards.
var FrenchCards = make([]rune, 52)

// TarotDeMarseilleSuits are the four suits present in a Tarot de Marseille deck of cards.
var TarotDeMarseilleSuits = []rune{'♣', '⚔', '⛾', '⛤'}

// TarotDeMarseilleValues are the fourteen values of cards for each suit in a Tarot de Marseille deck of cards.
var TarotDeMarseilleValues = []rune{'A', '2', '3', '4', '5', '6', '7', '8', '9', 'T', 'J', 'N', 'Q', 'K'}

// TarotDeMarseilleTrumps are the 22 trump cards in a Tarot de Marseille deck of cards.
var TarotDeMarseilleTrumps = []string{"0", "I", "II", "III", "IV", "V", "VI", "VII", "VIII", "IX", "X", "XI", "XII", "XIII", "XIV", "XV", "XVI", "XVII", "XVIII", "XIX", "XX", "XXI"}

// TarotDeMarseilleCards is a deck of 78 Tarot de Marseille cards.
var TarotDeMarseilleCards = make([]rune, 78)

func init() {
	card := 0
	for i := AceOfSpades; i < KingOfClubs; i++ {
		rank := i % 16
		// skip empty runes and jokers
		if rank >= 14 {
			continue
		}
		TarotDeMarseilleCards[card] = i
		// skip knights
		if rank == 11 {
			continue
		}
		FrenchCards[card] = i
		card++
	}
	for i := TheFool; i < TheWorld; i++ {
		TarotDeMarseilleCards[card] = i
		card++
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

// NewTarotDeMarseilleDeck builds a 78-card deck with four Italian suits
// (clubs, swords, cups, coins) of fourteen values (ace through king
// including knight between jack and queen) and 22 trumps (0 through XXI).
func NewTarotDeMarseilleDeck() *Deck {
	cards := make([]Card, 78)
	for i, c := range TarotDeMarseilleCards {
		cards[i] = Card(c)
	}
	return &Deck{cards: cards, draws: -1, tr: TranslateTarotDeMarseille}
}

// MaxDraws implements RandomObject interface.
func (d *Deck) MaxDraws() int {
	return len(d.cards)
}

// CountDistinctOutcomes implements RandomObject interface.
func (d *Deck) CountDistinctOutcomes(k int) *big.Int {
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

// NextOutcome implements RandomObject interface.
func (d *Deck) NextOutcome(k int) []rune {
	md := d.MaxDraws()
	if k > md {
		panic("k > MaxDraws")
	}
	if k < 0 {
		panic("k < 0")
	}
	if k == 0 {
		return nil
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

// TranslateFrench translates a playing card rune into a text name.
func TranslateFrench(r rune) (string, error) {
	if r < AceOfSpades || r > KingOfClubs {
		return "", fmt.Errorf("card '%c' is out of bounds", r)
	}
	suit := int(r-AceOfSpades) / 16
	if suit < 0 || suit >= len(FrenchSuits) {
		return "", fmt.Errorf("suit of card '%c' is unknown", r)
	}
	rank := int(r-AceOfSpades) % 16
	if rank < 0 || rank == 11 || rank > 13 {
		return "", fmt.Errorf("rank of card '%c' is unknown", r)
	}
	// skip knight
	if rank > 11 {
		rank -= 1
	}
	return string(FrenchValues[rank]) + string(FrenchSuits[suit]), nil
}

// TranslateTarotDeMarseille translates a playing card rune into a text name.
func TranslateTarotDeMarseille(r rune) (string, error) {
	if r < AceOfSpades || r > TheWorld {
		return "", fmt.Errorf("card '%c' is out of bounds", r)
	}
	// translate trumps first
	if r >= TheFool {
		rank := int(r - TheFool)
		return TarotDeMarseilleTrumps[rank], nil
	}

	suit := int(r-AceOfSpades) / 16
	if suit < 0 || suit >= len(TarotDeMarseilleSuits) {
		return "", fmt.Errorf("suit of card '%c' is unknown", r)
	}
	rank := int(r-AceOfSpades) % 16
	if rank < 0 || rank > 14 {
		return "", fmt.Errorf("rank of card '%c' is unknown", r)
	}
	return string(TarotDeMarseilleValues[rank]) + string(TarotDeMarseilleSuits[suit]), nil
}

// Translate implements RandomObject interface.
func (d *Deck) Translate(r rune) (string, error) {
	return d.tr(r)
}

// Card gets the nth card from the deck with no bounds checking
func (d *Deck) Card(n int) Card {
	return d.cards[n]
}
