package main

import (
	"fmt"
	"os"
)

func exitOnError(err error, message string) {
	if err != nil {
		fmt.Printf(message, err)
		os.Exit(1)
	}
}
