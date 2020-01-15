package cardware

import (
	"math/big"
	"reflect"
	"testing"
)

func TestDeck_MaxDraws(t *testing.T) {
	tests := []struct {
		name string
		d    *Deck
		want int
	}{
		{
			"zero",
			&Deck{cards: []Card{}},
			0,
		},
		{
			"one",
			&Deck{cards: []Card{'A'}},
			1,
		},
		{
			"standard",
			NewStandardFrenchDeck(),
			52,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.MaxDraws(); got != tt.want {
				t.Errorf("Deck.MaxDraws() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeck_CountDistinctOutcomes(t *testing.T) {

	fac52 := big.NewInt(int64(0))
	fac52.SetString("80658175170943878571660636856403766975289505440883277824000000000000", 10)

	type args struct {
		k int
	}
	tests := []struct {
		name      string
		d         *Deck
		args      args
		want      *big.Int
		wantPanic bool
	}{
		{
			name:      "zero",
			d:         &Deck{},
			args:      args{k: 0},
			want:      big.NewInt(int64(1)),
			wantPanic: false,
		},
		{
			name:      "one",
			d:         &Deck{cards: []Card{'A'}},
			args:      args{k: 1},
			want:      big.NewInt(int64(1)),
			wantPanic: false,
		},
		{
			name:      "standard-0",
			d:         NewStandardFrenchDeck(),
			args:      args{k: 0},
			want:      big.NewInt(int64(1)),
			wantPanic: false,
		},
		{
			name:      "standard-1",
			d:         NewStandardFrenchDeck(),
			args:      args{k: 1},
			want:      big.NewInt(int64(52)),
			wantPanic: false,
		},
		{
			name:      "standard-52",
			d:         NewStandardFrenchDeck(),
			args:      args{k: 52},
			want:      fac52,
			wantPanic: false,
		},
		{
			name:      "too-many",
			d:         &Deck{},
			args:      args{k: 1},
			want:      nil,
			wantPanic: true,
		},
		{
			name:      "negative",
			d:         &Deck{},
			args:      args{k: -1},
			want:      nil,
			wantPanic: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				r := recover()
				if (r != nil) != tt.wantPanic {
					t.Errorf("Deck.CountDistinctOutcomes() recover = %v, wantPanic = %v", r, tt.wantPanic)
				}
			}()
			if got := tt.d.CountDistinctOutcomes(tt.args.k); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Deck.CountDistinctOutcomes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeck_NextOutcome(t *testing.T) {
	type args struct {
		k int
	}
	tests := []struct {
		name      string
		d         *Deck
		args      args
		want      []rune
		wantPanic bool
	}{
		{
			name:      "zero",
			d:         &Deck{},
			args:      args{k: 0},
			want:      []rune{},
			wantPanic: false,
		},
		{
			name:      "one",
			d:         &Deck{cards: []Card{'A'}},
			args:      args{k: 1},
			want:      []rune{'A'},
			wantPanic: false,
		},
		{
			name:      "two",
			d:         &Deck{cards: []Card{'A', 'B'}},
			args:      args{k: 2},
			want:      []rune{'A', 'B'},
			wantPanic: false,
		},
		{
			name:      "two-1",
			d:         &Deck{cards: []Card{'A', 'B'}},
			args:      args{k: 1},
			want:      []rune{'A'},
			wantPanic: false,
		},
		{
			name:      "negative",
			d:         &Deck{cards: []Card{'A', 'B'}},
			args:      args{k: -1},
			want:      nil,
			wantPanic: true,
		},
		{
			name:      "too-many",
			d:         &Deck{cards: []Card{'A', 'B'}},
			args:      args{k: 3},
			want:      nil,
			wantPanic: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				r := recover()
				if (r != nil) != tt.wantPanic {
					t.Errorf("Deck.NextOutcome() recover = %v, wantPanic = %v", r, tt.wantPanic)
				}
			}()
			if got := tt.d.NextOutcome(tt.args.k); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Deck.NextOutcome() = %v, want %v", got, tt.want)
			}
		})
	}

	t.Run("standard-5-twice", func(t *testing.T) {
		want1 := []rune{'🂡', '🂢', '🂣', '🂤', '🂥'}
		want2 := []rune{'🂡', '🂢', '🂣', '🂥', '🂤'}
		deck := NewStandardFrenchDeck()
		if got := deck.NextOutcome(5); !reflect.DeepEqual(got, want1) {
			t.Errorf("Deck.NextOutcome() = %v, want %v", got, want1)
		}
		if got := deck.NextOutcome(5); !reflect.DeepEqual(got, want2) {
			t.Errorf("Deck.NextOutcome() = %v, want %v", got, want2)
		}
	})

	t.Run("count-3-3-twice", func(t *testing.T) {
		one := big.NewInt(int64(1))
		deck := &Deck{cards: []Card{'A', 'B', 'C'}}
		want := deck.CountDistinctOutcomes(3)
		i := big.NewInt(int64(0))
		for deck.NextOutcome(3) != nil {
			i.Add(i, one)
			// keep going!
		}
		if !reflect.DeepEqual(i, want) {
			t.Errorf("count = %v, want %v", i, want)
		}
		i = big.NewInt(int64(0))
		for deck.NextOutcome(3) != nil {
			i.Add(i, one)
			// keep going!
		}
		if !reflect.DeepEqual(i, want) {
			t.Errorf("count = %v, want %v", i, want)
		}
	})
}
