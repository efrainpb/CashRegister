package calculator

import "efrainpb/truefit-cashregister/internal/domain"

type GreedyChangeCalculator struct{}

func NewGreedyChangeCalculator() domain.ChangeCalculator {
	return &GreedyChangeCalculator{}
}

func (g *GreedyChangeCalculator) Calculate(changeInCents int, denominations []domain.Denomination) []domain.ChangeItem {
	items := make([]domain.ChangeItem, 0)

	for _, d := range denominations {
		if changeInCents <= 0 {
			break
		}

		count := changeInCents / d.Value
		if count == 0 {
			continue
		}
		items = append(items, domain.ChangeItem{Denomination: d, Count: count})
		changeInCents -= count * d.Value
	}

	return items
}
