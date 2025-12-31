package main

import (
	"context"
	"log"
	"os"

	"github.com/axelrhd/hagg/internal/ucli"
)

func main() {
	app := ucli.New()

	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
