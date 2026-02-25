package calculator

import (
	"efrainpb/truefit-cashregister/internal/domain"
	"testing"
)

var usdDenominations = []domain.Denomination{
	{Name: "Dollar", Value: 100},
	{Name: "Quarter", Value: 25},
	{Name: "Dime", Value: 10},
	{Name: "Nickel", Value: 5},
	{Name: "Penny", Value: 1},
}

func TestGreedyChangeCalculator_Calculate(t *testing.T) {
	calc := NewGreedyChangeCalculator()

	tests := []struct {
		name          string
		changeInCents int
		expected      []domain.ChangeItem
	}{
		{
			name:          "zero change",
			changeInCents: 0,
			expected:      []domain.ChangeItem{},
		},
		{
			name:          "exact dollar",
			changeInCents: 100,
			expected: []domain.ChangeItem{
				{Denomination: domain.Denomination{Name: "Dollar", Value: 100}, Count: 1},
			},
		},
		{
			name:          "41 cents",
			changeInCents: 41,
			expected: []domain.ChangeItem{
				{Denomination: domain.Denomination{Name: "Quarter", Value: 25}, Count: 1},
				{Denomination: domain.Denomination{Name: "Dime", Value: 10}, Count: 1},
				{Denomination: domain.Denomination{Name: "Nickel", Value: 5}, Count: 1},
				{Denomination: domain.Denomination{Name: "Penny", Value: 1}, Count: 1},
			},
		},
		{
			name:          "multiple dollars",
			changeInCents: 250,
			expected: []domain.ChangeItem{
				{Denomination: domain.Denomination{Name: "Dollar", Value: 100}, Count: 2},
				{Denomination: domain.Denomination{Name: "Quarter", Value: 25}, Count: 2},
			},
		},
		{
			name:          "only pennies",
			changeInCents: 3,
			expected: []domain.ChangeItem{
				{Denomination: domain.Denomination{Name: "Penny", Value: 1}, Count: 3},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calc.Calculate(tt.changeInCents, usdDenominations)
			if len(got) != len(tt.expected) {
				t.Fatalf("Calculate() returned %d items, want %d", len(got), len(tt.expected))
			}
			for i, item := range got {
				if item.Denomination.Name != tt.expected[i].Denomination.Name || item.Count != tt.expected[i].Count {
					t.Errorf("item[%d] = {%s, %d}, want {%s, %d}",
						i, item.Denomination.Name, item.Count,
						tt.expected[i].Denomination.Name, tt.expected[i].Count)
				}
			}
		})
	}
}

func TestGreedyChangeCalculator_TotalIsCorrect(t *testing.T) {
	calc := NewGreedyChangeCalculator()

	changeCases := []int{1, 5, 10, 25, 41, 99, 100, 167, 250, 999}
	for _, change := range changeCases {
		items := calc.Calculate(change, usdDenominations)
		total := 0
		for _, item := range items {
			total += item.Count * item.Denomination.Value
		}
		if total != change {
			t.Errorf("change=%d: items sum to %d", change, total)
		}
	}
}
