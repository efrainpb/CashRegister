package calculator

import (
	"efrainpb/truefit-cashregister/internal/domain"
	"math/rand"
	"time"
)

type RandomChangeCalculator struct{}

func NewRandomChangeCalculator() domain.ChangeCalculator {
	return &RandomChangeCalculator{}
}

func (r *RandomChangeCalculator) Calculate(changeInCents int, denominations []domain.Denomination) []domain.ChangeItem {
	shuffledDenominations := make([]domain.Denomination, len(denominations))
	copy(shuffledDenominations, denominations)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	rng.Shuffle(len(shuffledDenominations), func(i, j int) {
		shuffledDenominations[i], shuffledDenominations[j] = shuffledDenominations[j], shuffledDenominations[i]
	})

	items := make([]domain.ChangeItem, 0)
	remaining := changeInCents

	for _, d := range shuffledDenominations {
		if remaining <= 0 {
			break
		}
		count := remaining / d.Value
		if count == 0 {
			continue
		}
		items = append(items, domain.ChangeItem{Denomination: d, Count: count})
		remaining -= count * d.Value
	}

	return items
}
