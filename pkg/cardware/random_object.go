package cardware

import "math/big"

type RandomObject interface {
	MaxDraws() int
	CountDistinctOutcomes(k int) *big.Int
	NextOutcome(k int) []rune
	Translate(r rune) (string, error)
}
