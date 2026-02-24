package infrastructure

import (
	"context"
	"efrainpb/truefit-cashregister/internal/application"
	"efrainpb/truefit-cashregister/internal/domain"
	"efrainpb/truefit-cashregister/package/adapter/read"
	"efrainpb/truefit-cashregister/package/adapter/write"
	"fmt"
	"io"
	"os"
)

type FileController struct {
	processTransactions *application.ProcessTransactions
	source              io.Reader
}

func NewFileController(processTransactions *application.ProcessTransactions, source io.Reader) *FileController {
	return &FileController{
		processTransactions: processTransactions,
		source:              source,
	}
}

func (f *FileController) ProcessTransactions(currency domain.Currency) []domain.ChangeResult {
	reader := read.NewFileTransactionReader(f.source)
	transactions := reader.Read()
	if len(transactions) == 0 {
		return nil
	}

	results := f.processTransactions.Process(transactions, currency.Denominations())

	// output to console
	consoleWriter := write.NewConsoleChangeWriter(results, os.Stdout, currency)
	if err := consoleWriter.Writer(context.Background()); err != nil {
		fmt.Println("error writing to console:", err)
	}

	// output to file
	/*fileWriter := write.NewFileChangeWriter("output.txt", results, os.Stdout, currency)
	if err := fileWriter.Writer(context.Background()); err != nil {
		fmt.Println("error writing to file:", err)
	}*/

	return results
}
