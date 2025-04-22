package main

import "fmt"

func findIndex[T comparable](s []T, x T) int {

	for index, value := range s {
		if value == x {
			return index
		}
	}
	return -1
}

func main() {
	numbers := []int{1, 2, 3, 4, 5}
	fmt.Println(findIndex(numbers, 4))

	strings := []string{"a", "b", "c"}
	fmt.Println(findIndex(strings, "d"))
}

// generic types
type List[T any] struct {
	next *List[T]
	val  T
}
