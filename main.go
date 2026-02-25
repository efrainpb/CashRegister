package main

import (
	"efrainpb/truefit-cashregister/internal"
	"fmt"
	"os"
)

func main() {
	if err := app.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
