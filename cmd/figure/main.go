package main

import (
	"log"
	"pureProject/internal/figure"
)

func main() {
	if err := figure.NewApp().Run(); err != nil {
		log.Fatal(err)
	}
}
