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
		name := pluralize(item.Denomination.Name, item.Count)
		parts = append(parts, fmt.Sprintf("%d %s", item.Count, name))
	}
	return strings.Join(parts, ", ")
}

func pluralize(name string, count int) string {
	if count == 1 {
		return name
	}
	if strings.HasSuffix(name, "nny") {
		return name[:len(name)-1] + "ies"
	}
	return name + "s"
}
