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

	fmt.Println("==================")

	queue := structs.Queue{}

	for i := 0; i < 15; i++ {
		queue.Push(i)
	}

	fmt.Println(queue)
	fmt.Println(queue.Pop())
	fmt.Println(queue.Pop())
	fmt.Println(queue.Pop())
	fmt.Println(queue)

}
