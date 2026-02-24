package app

import (
	"efrainpb/truefit-cashregister/internal/application"
	"efrainpb/truefit-cashregister/internal/domain"
	"efrainpb/truefit-cashregister/internal/infrastructure"
	"efrainpb/truefit-cashregister/package/currency"
	"fmt"
	"os"
)

func Run() error {
	f, err := os.Open("input.txt")
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer f.Close()

	currrency := currency.NewUSD()

	processTransactions := application.NewProcessTransactions(domain.AmountDivisor)
	controller := infrastructure.NewFileController(processTransactions, f)
	controller.ProcessTransactions(currrency)
	return nil
}
