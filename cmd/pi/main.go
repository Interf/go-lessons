package main

import (
	"fmt"
	"go-lessons/internal/pi"
	"os"
)

func main() {

	if err := pi.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
