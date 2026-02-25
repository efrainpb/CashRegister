package currency

import (
	"efrainpb/truefit-cashregister/internal/domain"
	"testing"
)

func TestUSD_Denominations(t *testing.T) {
	usd := NewUSD()
	denoms := usd.Denominations()

	expected := []struct {
		name  string
		value int
	}{
		{"Dollar", 100},
		{"Quarter", 25},
		{"Dime", 10},
		{"Nickel", 5},
		{"Penny", 1},
	}

	if len(denoms) != len(expected) {
		t.Fatalf("Denominations() returned %d items, want %d", len(denoms), len(expected))
	}

	for i, d := range denoms {
		if d.Name != expected[i].name || d.Value != expected[i].value {
			t.Errorf("denomination[%d] = {%s, %d}, want {%s, %d}",
				i, d.Name, d.Value, expected[i].name, expected[i].value)
		}
	}
}

func TestUSD_DenominationsDescendingOrder(t *testing.T) {
	usd := NewUSD()
	denoms := usd.Denominations()

	for i := 1; i < len(denoms); i++ {
		if denoms[i].Value >= denoms[i-1].Value {
			t.Errorf("denominations not in descending order: %d >= %d", denoms[i].Value, denoms[i-1].Value)
		}
	}
}

func TestUSD_FormatChange(t *testing.T) {
	usd := NewUSD()

	tests := []struct {
		name     string
		items    []domain.ChangeItem
		expected string
	}{
		{
			name:     "empty",
			items:    []domain.ChangeItem{},
			expected: "",
		},
		{
			name: "single dollar",
			items: []domain.ChangeItem{
				{Denomination: domain.Denomination{Name: "Dollar", Value: 100}, Count: 1},
			},
			expected: "1 Dollar",
		},
		{
			name: "multiple dollars",
			items: []domain.ChangeItem{
				{Denomination: domain.Denomination{Name: "Dollar", Value: 100}, Count: 2},
			},
			expected: "2 Dollars",
		},
		{
			name: "single penny",
			items: []domain.ChangeItem{
				{Denomination: domain.Denomination{Name: "Penny", Value: 1}, Count: 1},
			},
			expected: "1 Penny",
		},
		{
			name: "multiple pennies",
			items: []domain.ChangeItem{
				{Denomination: domain.Denomination{Name: "Penny", Value: 1}, Count: 3},
			},
			expected: "3 Pennies",
		},
		{
			name: "mixed denominations",
			items: []domain.ChangeItem{
				{Denomination: domain.Denomination{Name: "Dollar", Value: 100}, Count: 1},
				{Denomination: domain.Denomination{Name: "Quarter", Value: 25}, Count: 2},
				{Denomination: domain.Denomination{Name: "Penny", Value: 1}, Count: 4},
			},
			expected: "1 Dollar, 2 Quarters, 4 Pennies",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := usd.FormatChange(tt.items)
			if got != tt.expected {
				t.Errorf("FormatChange() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestPluralize(t *testing.T) {
	cases := []struct {
		input    string
		count    int
		expected string
	}{
		{"Dollar", 1, "Dollar"},
		{"Dollar", 2, "Dollars"},
		{"Penny", 1, "Penny"},
		{"Penny", 2, "Pennies"},
		{"Penny", 5, "Pennies"},
		{"Quarter", 1, "Quarter"},
		{"Quarter", 3, "Quarters"},
		{"Nickel", 1, "Nickel"},
		{"Nickel", 5, "Nickels"},
		{"Dime", 1, "Dime"},
		{"Dime", 10, "Dimes"},
	}
	for _, c := range cases {
		got := pluralize(c.input, c.count)
		if got != c.expected {
			t.Errorf("pluralize(%q, %d) = %q, want %q", c.input, c.count, got, c.expected)
		}
	}
}
