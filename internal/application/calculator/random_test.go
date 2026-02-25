package calculator

import (
	"testing"
)

func TestRandomChangeCalculator_TotalIsCorrect(t *testing.T) {
	calc := NewRandomChangeCalculator()

	changeCases := []int{1, 5, 10, 25, 41, 99, 100, 167, 250, 999}
	for _, change := range changeCases {
		// Run multiple times to account for randomness
		for i := 0; i < 10; i++ {
			items := calc.Calculate(change, usdDenominations)
			total := 0
			for _, item := range items {
				total += item.Count * item.Denomination.Value
			}
			if total != change {
				t.Errorf("change=%d (run %d): items sum to %d", change, i+1, total)
			}
		}
	}
}

func TestRandomChangeCalculator_ZeroChange(t *testing.T) {
	calc := NewRandomChangeCalculator()
	items := calc.Calculate(0, usdDenominations)
	if len(items) != 0 {
		t.Errorf("expected 0 items for zero change, got %d", len(items))
	}
}

func TestRandomChangeCalculator_DoesNotMutateInput(t *testing.T) {
	calc := NewRandomChangeCalculator()

	original := make([]string, len(usdDenominations))
	for i, d := range usdDenominations {
		original[i] = d.Name
	}

	calc.Calculate(100, usdDenominations)

	for i, d := range usdDenominations {
		if d.Name != original[i] {
			t.Errorf("denomination[%d] was mutated: got %s, want %s", i, d.Name, original[i])
		}
	}
}
