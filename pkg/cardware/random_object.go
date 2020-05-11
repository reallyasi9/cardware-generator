package cardware

import "math/big"

// RandomObject is an interface to a thing that can draw random elements from a fixed set
// (like a deck of cards or a single die).
type RandomObject interface {
	MaxDraws() int
	CountDistinctOutcomes(k int) *big.Int
	NextOutcome(k int) []rune
	Translate(r rune) (string, error)
}
