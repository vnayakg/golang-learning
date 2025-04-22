package filter

func filterEvenNumbers(numbers []int) []int {
	return filterNumbersOnPredicate(numbers, func(num int) bool { return num%2 == 0 })
}

func filterOddNumbers(numbers []int) []int {
	return filterNumbersOnPredicate(numbers, func(num int) bool { return num%2 != 0 })
}

func filterNumbersOnPredicate(numbers []int, predicate func(int) bool) []int {
	filteredNumbers := make([]int, 0, len(numbers))
	for _, num := range numbers {
		if predicate(num) {
			filteredNumbers = append(filteredNumbers, num)
		}
	}
	return filteredNumbers
}
