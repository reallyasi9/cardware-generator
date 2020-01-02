package cardware

import "math/big"

type RandomObject interface {
	MaxDraws() int
	NumPermutations(k int) *big.Int
	NumCombinations(k int) *big.Int
	NextPermutation(k int) []rune
	NextCombination(k int) []rune
	Translate(r rune) (string, error)
}
