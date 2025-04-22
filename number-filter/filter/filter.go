package filter

func filterEvenNumbers(numbers []int) []int {
	evenNumbers := make([]int, 0, len(numbers))
	for _, num := range numbers {
		if num % 2 == 0 {
			evenNumbers = append(evenNumbers, num)
		}
	}
	return evenNumbers
}

func filterOddNumbers(numbers []int) []int {
	oddNumbers := make([]int, 0, len(numbers))
	for _, num := range numbers {
		if num % 2 != 0 {
			oddNumbers = append(oddNumbers, num)
		}
	}
	return oddNumbers
}
