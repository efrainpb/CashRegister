package currency

import (
	"efrainpb/truefit-cashregister/internal/domain"
	"fmt"
	"strings"
)

type USD struct{}

func NewUSD() *USD {
	return &USD{}
}

func (u *USD) Denominations() []domain.Denomination {
	return []domain.Denomination{
		{Name: "Dollar", Value: 100},
		{Name: "Quarter", Value: 25},
		{Name: "Dime", Value: 10},
		{Name: "Nickel", Value: 5},
		{Name: "Penny", Value: 1},
	}
}

func (u *USD) FormatChange(items []domain.ChangeItem) string {
	parts := make([]string, 0)
	for _, item := range items {
		parts = append(parts, fmt.Sprintf("%d %s", item.Count, item.Denomination.Name))
	}
	return strings.Join(parts, ", ")
}
