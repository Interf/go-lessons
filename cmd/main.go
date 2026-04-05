package main

import (
	"fmt"

	"example.com/m/v2/internal/structs"
)

func main() {

	stack := structs.Stack{}

	for i := 0; i < 15; i++ {
		stack.Push(structs.Item{
			Value: i,
		})
	}

	fmt.Println(stack)
	fmt.Println(stack.Pop())
	fmt.Println(stack.Pop())
	fmt.Println(stack.Pop())
	fmt.Println(stack)

}
