package read

import (
	"bufio"
	"efrainpb/truefit-cashregister/internal/domain"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
)

type FileTransactionReader struct {
	source io.Reader
}

func NewFileTransactionReader(source io.Reader) *FileTransactionReader {
	return &FileTransactionReader{source: source}
}

func (f *FileTransactionReader) Read() []domain.Transaction {
	var transactions []domain.Transaction

	scanner := bufio.NewScanner(f.source)
	lineNumber := 0
	for scanner.Scan() {
		lineNumber++
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		tx, err := parseLine(line)
		if err != nil {
			fmt.Println("line", lineNumber, "is invalid:", err)
			continue
		}
		transactions = append(transactions, tx)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("error reading file:", err)
	}
	return transactions
}

func parseLine(line string) (domain.Transaction, error) {
	parts := strings.SplitN(line, ",", 2)
	if len(parts) != 2 {
		return domain.Transaction{}, fmt.Errorf("invalid line: %s", line)
	}

	owed, err := parseCents(strings.TrimSpace(parts[0]))
	if err != nil {
		return domain.Transaction{}, fmt.Errorf("invalid amount owed: %s", line)
	}

	paid, err := parseCents(strings.TrimSpace(parts[1]))
	if err != nil {
		return domain.Transaction{}, fmt.Errorf("invalid amount paid: %s", line)
	}

	if paid < owed {
		return domain.Transaction{}, fmt.Errorf("invalid amount paid: %s", line)
	}

	return domain.Transaction{
		AmountOwed: owed,
		AmountPaid: paid,
	}, nil
}

func parseCents(amount string) (int, error) {
	val, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid amount: %s", amount)
	}
	return int(math.Round(val * 100)), nil
}
