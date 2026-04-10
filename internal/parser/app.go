package parser

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"parser/internal/parser/pack"
	"parser/internal/parser/unpack"
)

func Run() {
	var packMode bool
	var unpackMode bool
	var input string
	var isDaemonMode bool

	flag.BoolVar(&packMode, "pack", false, "pack mode")
	flag.BoolVar(&unpackMode, "unpack", false, "unpack mode")
	flag.StringVar(&input, "input", "", "input string")
	flag.BoolVar(&isDaemonMode, "daemon", false, "daemon mode")

	flag.Parse()

	var result string
	var err error

	if packMode {
		result = pack.Pack(input)
	} else if unpackMode {

		if isDaemonMode {
			scanner := bufio.NewScanner(os.Stdin)

			for scanner.Scan() {
				text := scanner.Text()

				result, err = unpack.Unpack(text)
				if err != nil {
					fmt.Println("Error:", err)
				} else {
					fmt.Println(result)
				}
			}
		} else {
			result, err = unpack.Unpack(input)
		}

	}

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
}
