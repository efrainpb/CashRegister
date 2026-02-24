package domain

type Denomination struct {
	Name  string
	Value int
}

type Transaction struct {
	AmountOwed int
	AmountPaid int
}

func (t Transaction) Change() int {
	return t.AmountPaid - t.AmountOwed
}

type ChangeItem struct {
	Denomination Denomination
	Count        int
}

type ChangeResult struct {
	Transaction Transaction
	Items       []ChangeItem
}

type Currency interface {
	Denominations() []Denomination
	FormatChange(items []ChangeItem) string
}

type ChangeCalculator interface {
	Calculate(changeInCents int, denominations []Denomination) []ChangeItem
}
