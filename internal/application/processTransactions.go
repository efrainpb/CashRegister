package application

import (
	"efrainpb/truefit-cashregister/internal/application/calculator"
	"efrainpb/truefit-cashregister/internal/domain"
)

func NewProcessTransactions(
	amountDivisor int,
) *ProcessTransactions {
	return &ProcessTransactions{
		amountDivisor: amountDivisor,
	}
}

type ProcessTransactions struct {
	amountDivisor int
}

func (p *ProcessTransactions) Process(transactions []domain.Transaction, denominations []domain.Denomination) []domain.ChangeResult {
	results := make([]domain.ChangeResult, 0)

	if len(transactions) == 0 {
		return nil
	}

	for _, t := range transactions {
		results = append(results, p.processTransaction(t, denominations))
	}
	return results
}

func (p *ProcessTransactions) processTransaction(t domain.Transaction, denominations []domain.Denomination) domain.ChangeResult {
	if t.Change() < 0 {
		return domain.ChangeResult{}
	}

	var changeResult domain.ChangeResult
	var strategy domain.ChangeCalculator

	if t.Change()%p.amountDivisor == 0 {
		strategy = calculator.NewRandomChangeCalculator()
	} else {
		strategy = calculator.NewGreedyChangeCalculator()
	}

	changeResult.Items = strategy.Calculate(t.Change(), denominations)
	changeResult.Transaction = t
	return changeResult
}
