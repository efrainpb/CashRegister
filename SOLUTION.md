# Solution Document — Cash Register

## Overview

This application calculates change for cashier transactions. It reads a flat file of transactions, computes the optimal denomination breakdown for each, and outputs the results to the console (and optionally to a file). The key business rule is that when the change amount is evenly divisible by 3, denominations are selected randomly instead of greedily.

---

## Architecture

The solution follows **Hexagonal Architecture** (Ports & Adapters), separating the domain logic from I/O concerns. This makes the core logic independently testable and allows swapping input/output mechanisms without touching business rules.

```
main.go
└── internal/
    ├── app.go                          # Composition root — wires everything together
    ├── domain/
    │   └── models.go                   # Core types and interfaces (Currency, ChangeCalculator)
    ├── application/
    │   ├── processTransactions.go      # Orchestration: strategy selection per transaction
    │   └── calculator/
    │       ├── greedy.go               # Greedy algorithm (minimum coins)
    │       └── random.go               # Random denomination algorithm
    └── infrastructure/
        └── fileController.go           # Drives the use case from a file source
package/
    ├── adapter/
    │   ├── read/
    │   │   └── file.go                 # Parses input.txt into domain.Transaction slice
    │   └── write/
    │       ├── console.go              # Writes formatted output to stdout
    │       └── file.go                 # Writes formatted output to output.txt
    └── currency/
        └── usd.go                      # USD denomination list and change formatting
```

---

## Data Flow

```
input.txt
    │
    ▼
FileTransactionReader.Read()
    │  Parses each line ("2.13,3.00") into domain.Transaction{AmountOwed, AmountPaid}
    │  Amounts are stored in cents (integer) to avoid float precision errors
    ▼
ProcessTransactions.Process()
    │  For each transaction, computes Change = AmountPaid - AmountOwed
    │  Selects a strategy:
    │    • Change % 3 == 0  →  RandomChangeCalculator
    │    • Otherwise        →  GreedyChangeCalculator
    ▼
ChangeCalculator.Calculate()
    │  Returns []ChangeItem{Denomination, Count}
    ▼
ConsoleChangeWriter / FileChangeWriter
    │  Calls currency.FormatChange() to build the output string
    │  e.g. "3 quarters, 1 dime, 3 pennies"
    ▼
stdout  (and optionally output.txt)
```

---

## Key Design Decisions

### 1. Integer arithmetic for money

Dollar amounts are converted to **cents** (`int`) immediately on parsing, using `math.Round(val * 100)`. This avoids floating-point rounding errors that are common when doing arithmetic on currency values represented as `float64`.

```go
// package/adapter/read/file.go:73
func parseCents(amount string) (int, error) {
    val, err := strconv.ParseFloat(amount, 64)
    ...
    return int(math.Round(val * 100)), nil
}
```

### 2. Strategy Pattern for change calculation

The `domain.ChangeCalculator` interface is the central extension point:

```go
// internal/domain/models.go:34
type ChangeCalculator interface {
    Calculate(changeInCents int, denominations []Denomination) []ChangeItem
}
```

`ProcessTransactions` selects the concrete implementation at runtime based on the divisibility rule. Adding a new calculation strategy (e.g., "fewest bills only") requires only implementing this interface — no changes to orchestration or I/O code.

### 3. Greedy algorithm

When change is **not** divisible by 3, the greedy algorithm iterates denominations from largest to smallest and takes as many of each as possible:

```
Change = 88 cents
→ 3 quarters  (75 cents)  remaining = 13
→ 1 dime      (10 cents)  remaining = 3
→ 3 pennies   ( 3 cents)  remaining = 0
→ output: "3 Quarters, 1 Dime, 3 Pennies"
```

### 4. Random algorithm

When owed amount **is** divisible by 3, the denomination list is shuffled with a time-seeded RNG, then the same greedy logic is applied to the shuffled order. The math is always correct; only the denomination mix varies between runs.

```go
// internal/application/calculator/random.go
rng := rand.New(rand.NewSource(time.Now().UnixNano()))
rng.Shuffle(len(shuffledDenominations), func(i, j int) { ... })
```

### 5. Configurable divisor

The divisibility threshold (`3`) is defined as a named constant in the domain layer and injected into `ProcessTransactions` at construction time:

```go
// internal/domain/models.go
const AmountDivisor = 3

// internal/app.go
processTransactions := application.NewProcessTransactions(domain.AmountDivisor)
```

Changing the business rule from "divisible by 3" to any other value requires updating a single constant.

### 6. Currency as a plug-in

The `domain.Currency` interface decouples denomination definitions and formatting from the core logic:

```go
type Currency interface {
    Denominations() []Denomination
    FormatChange(items []ChangeItem) string
}
```

The only concrete implementation today is `USD` (`package/currency/usd.go`). Supporting a new locale (e.g., euros for a French client) means adding a new `Currency` implementation — zero changes to the domain or application layers.

---

## USD Denominations

| Name    | Value (cents) |
|---------|--------------|
| Dollar  | 100          |
| Quarter | 25           |
| Dime    | 10           |
| Nickel  | 5            |
| Penny   | 1            |

Pluralization follows English rules: `Penny → Pennies`, all others get an `s` suffix.

---

## Input / Output

**Input file** (`input.txt`): one transaction per line, comma-separated.

```
2.12,3.00
1.97,2.00
3.33,5.00
```

**Output** (console / `output.txt`):

```
3 Quarters, 1 Dime, 3 Pennies
3 Pennies
1 Dollar, 1 Quarter, 6 Nickels, 12 Pennies   ← random (divisible by 3)
```

Invalid lines (malformed, unparseable amounts, or paid < owed) are logged and skipped; processing continues for the remaining lines.

---

## Running Tests

### Run all tests

```bash
go test ./...
```

### Run with coverage report

```bash
go test ./... -coverprofile=coverage.out
go tool cover -func=coverage.out
```

### View coverage in browser

```bash
go tool cover -html=coverage.out
```

### Run a specific package

```bash
go test ./internal/application/calculator/...
```

---

## Running the Application

### With Go

```bash
go run .
```

### With Docker

```bash
docker build -t cashregister .
docker run --rm -v $(pwd):/data cashregister
```

The Docker image mounts the host directory at `/data`, reading `input.txt` and writing `output.txt` there.

---

## Extension Points

| Requirement change | Where to touch |
|--------------------|----------------|
| Change the random divisor | `internal/domain/models.go` → `AmountDivisor` |
| Add a new calculation strategy | Implement `domain.ChangeCalculator`, add selection logic in `processTransactions.go` |
| Support a new currency/locale | Implement `domain.Currency`, pass it in `internal/app.go` |
| Add a new output target (DB, API) | Add a new writer under `package/adapter/write/` |
| Add a new input source (stdin, HTTP) | Add a new reader under `package/adapter/read/` |
