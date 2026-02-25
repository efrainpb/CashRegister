package domain

import "testing"

func TestTransaction_Change(t *testing.T) {
	tests := []struct {
		name     string
		tx       Transaction
		expected int
	}{
		{"exact change", Transaction{AmountOwed: 100, AmountPaid: 100}, 0},
		{"positive change", Transaction{AmountOwed: 75, AmountPaid: 200}, 125},
		{"negative change (overpaid by owed)", Transaction{AmountOwed: 200, AmountPaid: 100}, -100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tx.Change(); got != tt.expected {
				t.Errorf("Change() = %d, want %d", got, tt.expected)
			}
		})
	}
}
