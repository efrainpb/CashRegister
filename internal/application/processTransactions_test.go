package application

import (
	"efrainpb/truefit-cashregister/internal/domain"
	"testing"
)

var denominations = []domain.Denomination{
	{Name: "Dollar", Value: 100},
	{Name: "Quarter", Value: 25},
	{Name: "Dime", Value: 10},
	{Name: "Nickel", Value: 5},
	{Name: "Penny", Value: 1},
}

func TestProcessTransactions_EmptyInput(t *testing.T) {
	p := NewProcessTransactions(domain.AmountDivisor)
	result := p.Process([]domain.Transaction{}, denominations)
	if result != nil {
		t.Errorf("expected nil for empty transactions, got %v", result)
	}
}

func TestProcessTransactions_NilInput(t *testing.T) {
	p := NewProcessTransactions(domain.AmountDivisor)
	result := p.Process(nil, denominations)
	if result != nil {
		t.Errorf("expected nil for nil transactions, got %v", result)
	}
}

func TestProcessTransactions_ChangeIsCorrect(t *testing.T) {
	p := NewProcessTransactions(domain.AmountDivisor)

	transactions := []domain.Transaction{
		{AmountOwed: 175, AmountPaid: 200}, // change = 25 cents, owed not divisible by 3 → greedy
		{AmountOwed: 150, AmountPaid: 300}, // change = 150 cents, owed divisible by 3 → random
	}

	results := p.Process(transactions, denominations)

	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}

	for i, r := range results {
		total := 0
		for _, item := range r.Items {
			total += item.Count * item.Denomination.Value
		}
		expected := transactions[i].Change()
		if total != expected {
			t.Errorf("transaction[%d]: items sum to %d, want %d", i, total, expected)
		}
	}
}

func TestProcessTransactions_TransactionPreserved(t *testing.T) {
	p := NewProcessTransactions(domain.AmountDivisor)

	tx := domain.Transaction{AmountOwed: 100, AmountPaid: 200}
	results := p.Process([]domain.Transaction{tx}, denominations)

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Transaction != tx {
		t.Errorf("transaction not preserved: got %+v, want %+v", results[0].Transaction, tx)
	}
}

func TestProcessTransactions_NegativeChange(t *testing.T) {
	p := NewProcessTransactions(domain.AmountDivisor)

	// AmountPaid < AmountOwed → change is negative → empty ChangeResult
	tx := domain.Transaction{AmountOwed: 200, AmountPaid: 100}
	results := p.Process([]domain.Transaction{tx}, denominations)

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Items != nil || results[0].Transaction != (domain.Transaction{}) {
		t.Errorf("expected empty ChangeResult for negative change, got %+v", results[0])
	}
}

func TestProcessTransactions_DivisorSelectsStrategy(t *testing.T) {
	p := NewProcessTransactions(3)

	// AmountOwed = 300 is divisible by 3 → random strategy
	// AmountOwed = 100 is not divisible by 3 → greedy strategy
	// We can't assert which strategy was used, but we can assert the total is correct.
	transactions := []domain.Transaction{
		{AmountOwed: 300, AmountPaid: 400},
		{AmountOwed: 100, AmountPaid: 200},
	}

	results := p.Process(transactions, denominations)
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}

	for i, r := range results {
		total := 0
		for _, item := range r.Items {
			total += item.Count * item.Denomination.Value
		}
		expected := transactions[i].Change()
		if total != expected {
			t.Errorf("transaction[%d]: items sum to %d, want %d", i, total, expected)
		}
	}
}
