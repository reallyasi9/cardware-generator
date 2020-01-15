package cardware

import "math/big"

// Combined embeds a DiceBag and a Deck.
type Combined struct {
	RandomObject
	DiceBag
	Deck
	lastDice []rune
}

// NewCombined creates a new DiceBag and Deck simultaneously
func NewCombined(dice []int) *Combined {
	db := NewDiceBag(dice)
	deck := NewStandardFrenchDeck()
	return &Combined{DiceBag: *db, Deck: *deck}
}

func (c *Combined) splitDraws(k int) (int, int) {
	kdice := k
	kcards := 0
	if k > c.DiceBag.MaxDraws() {
		kdice = c.DiceBag.MaxDraws()
		kcards = k - kdice
	}
	return kdice, kcards
}

// MaxDraws implements RandomObject interface.
func (c *Combined) MaxDraws() int {
	return c.DiceBag.MaxDraws() * c.Deck.MaxDraws()
}

// CountDistinctOutcomes implements RandomObject interface.
func (c *Combined) CountDistinctOutcomes(k int) *big.Int {
	kdice, kcards := c.splitDraws(k)
	out := c.DiceBag.CountDistinctOutcomes(kdice)
	return out.Mul(out, c.Deck.CountDistinctOutcomes(kcards))
}

// NextOutcome implements RandomObject interface.
func (c *Combined) NextOutcome(k int) []rune {
	kdice, kcards := c.splitDraws(k)
	if c.lastDice == nil {
		c.lastDice = c.DiceBag.NextOutcome(kdice)
	}
	cards := c.Deck.NextOutcome(kcards)
	if cards == nil {
		c.lastDice = c.DiceBag.NextOutcome(kdice)
		cards = c.Deck.NextOutcome(kcards)
	}
	if c.lastDice == nil {
		return nil // nothing left to draw.
	}
	return append(c.lastDice, cards...)
}

// Translate implements RandomObject interface.
func (c *Combined) Translate(r rune) (string, error) {
	if r < AceOfSpades {
		return c.DiceBag.Translate(r)
	}
	return c.Deck.Translate(r)
}
