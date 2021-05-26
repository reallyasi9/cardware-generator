package cardware

import (
	"math/big"
	"math/rand"
)

// RandomObject is an interface to a thing that can draw random elements from a fixed set
// (like a deck of cards or a single die).
type RandomObject interface {
	MaxDraws() int
	CountDistinctOutcomes(k int) *big.Int
	NextOutcome(k int) []rune
	Translate(r rune) (string, error)
}

// Shuffler is an interface to a thing that can be shuffled.
type Shuffler interface {
	Shuffle(src rand.Source)
}

// Drawer is an interface to a collection of randomly-sorted objects, from which a single value can be drawn.
type Drawer interface {
	HasRemaining() bool
	Draw() (string, error)
}
