package read

import (
	"strings"
	"testing"
)

func TestFileTransactionReader_Read(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantLen  int
		wantOwed []int
		wantPaid []int
	}{
		{
			name:     "valid single transaction",
			input:    "1.75,2.00\n",
			wantLen:  1,
			wantOwed: []int{175},
			wantPaid: []int{200},
		},
		{
			name:     "multiple valid transactions",
			input:    "1.75,2.00\n3.00,5.00\n0.50,1.00\n",
			wantLen:  3,
			wantOwed: []int{175, 300, 50},
			wantPaid: []int{200, 500, 100},
		},
		{
			name:    "empty input",
			input:   "",
			wantLen: 0,
		},
		{
			name:     "blank lines skipped",
			input:    "\n1.75,2.00\n\n",
			wantLen:  1,
			wantOwed: []int{175},
			wantPaid: []int{200},
		},
		{
			name:     "invalid line skipped",
			input:    "bad_line\n1.75,2.00\n",
			wantLen:  1,
			wantOwed: []int{175},
			wantPaid: []int{200},
		},
		{
			name:     "paid less than owed skipped",
			input:    "5.00,2.00\n1.75,2.00\n",
			wantLen:  1,
			wantOwed: []int{175},
			wantPaid: []int{200},
		},
		{
			name:     "exact payment (zero change)",
			input:    "2.00,2.00\n",
			wantLen:  1,
			wantOwed: []int{200},
			wantPaid: []int{200},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewFileTransactionReader(strings.NewReader(tt.input))
			txs := r.Read()

			if len(txs) != tt.wantLen {
				t.Fatalf("Read() returned %d transactions, want %d", len(txs), tt.wantLen)
			}

			for i, tx := range txs {
				if tx.AmountOwed != tt.wantOwed[i] {
					t.Errorf("txs[%d].AmountOwed = %d, want %d", i, tx.AmountOwed, tt.wantOwed[i])
				}
				if tx.AmountPaid != tt.wantPaid[i] {
					t.Errorf("txs[%d].AmountPaid = %d, want %d", i, tx.AmountPaid, tt.wantPaid[i])
				}
			}
		})
	}
}

func TestParseLine(t *testing.T) {
	tests := []struct {
		name      string
		line      string
		wantOwed  int
		wantPaid  int
		wantError bool
	}{
		{"valid", "1.75,2.00", 175, 200, false},
		{"whole numbers", "3,5", 300, 500, false},
		{"exact change", "1.00,1.00", 100, 100, false},
		{"no comma", "1.75", 0, 0, true},
		{"invalid owed", "abc,2.00", 0, 0, true},
		{"invalid paid", "1.75,abc", 0, 0, true},
		{"paid less than owed", "5.00,2.00", 0, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx, err := parseLine(tt.line)
			if tt.wantError {
				if err == nil {
					t.Errorf("parseLine(%q) expected error, got nil", tt.line)
				}
				return
			}
			if err != nil {
				t.Fatalf("parseLine(%q) unexpected error: %v", tt.line, err)
			}
			if tx.AmountOwed != tt.wantOwed {
				t.Errorf("AmountOwed = %d, want %d", tx.AmountOwed, tt.wantOwed)
			}
			if tx.AmountPaid != tt.wantPaid {
				t.Errorf("AmountPaid = %d, want %d", tx.AmountPaid, tt.wantPaid)
			}
		})
	}
}

func TestParseCents(t *testing.T) {
	tests := []struct {
		input    string
		expected int
		wantErr  bool
	}{
		{"1.75", 175, false},
		{"2.00", 200, false},
		{"0.01", 1, false},
		{"0.10", 10, false},
		{"10.00", 1000, false},
		{"abc", 0, true},
		{"", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := parseCents(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("parseCents(%q) expected error, got nil", tt.input)
				}
				return
			}
			if err != nil {
				t.Fatalf("parseCents(%q) unexpected error: %v", tt.input, err)
			}
			if got != tt.expected {
				t.Errorf("parseCents(%q) = %d, want %d", tt.input, got, tt.expected)
			}
		})
	}
}
