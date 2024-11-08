package main

import (
	"log"

	"github.com/nkbhasker/go-pgx-transaction-example/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
