package cardware

import (
	"fmt"
	"math/big"

	"gonum.org/v1/gonum/stat/combin"
)

// DiceBag represents a bag of individual dice.
type DiceBag struct {
	RandomObject
	dice  []int
	using int
	cg    *combin.CartesianGenerator
}

// NewDiceBag creates a new bag of dice from a collection of dice.
func NewDiceBag(dice []int) *DiceBag {
	d := make([]int, len(dice))
	copy(d, dice)
	return &DiceBag{dice: d}
}

// MaxDraws implements RandomObject interface.
func (d *DiceBag) MaxDraws() int {
	return len(d.dice)
}

// CountDistinctOutcomes implements RandomObject interface.
func (d *DiceBag) CountDistinctOutcomes(k int) *big.Int {
	md := d.MaxDraws()
	if k > md {
		panic("k > MaxDraws")
	}
	if k < 0 {
		panic("k < 0")
	}
	n := big.NewInt(1)
	if k == 0 {
		return n
	}
	for i := 0; i < k; i++ {
		n.Mul(n, big.NewInt(int64(d.dice[i])))
	}
	return n
}

// NextOutcome implements RandomObject interface.
func (d *DiceBag) NextOutcome(k int) []rune {
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
	if k != d.using || d.cg == nil {
		d.using = k
		d.cg = combin.NewCartesianGenerator(d.dice[:k])
	}
	if !d.cg.Next() {
		d.cg = nil
		return nil
	}
	p := make([]int, k)
	d.cg.Product(p)
	out := make([]rune, k)
	for i, x := range p {
		out[i] = rune(x)
	}
	return out
}

// Translate implements RandomObject interface.
func (d *DiceBag) Translate(r rune) (string, error) {
	return fmt.Sprintf("[%d]", int(r)+1), nil
}
