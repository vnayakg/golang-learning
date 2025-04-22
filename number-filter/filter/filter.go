package filter

func filterEvenNumbers(numbers []int) []int {
	return filterNumbersOnPredicate(numbers, isEven)
}

func filterOddNumbers(numbers []int) []int {
	return filterNumbersOnPredicate(numbers, isOdd)
}

func filterPrimeNumbers(numbers []int) []int {
	return filterNumbersOnPredicate(numbers, isPrime)
}

func filterOddPrimeNumbers(number []int) []int {
	return filterNumbersOnPredicate(number, func(num int) bool { return isOdd(num) && isPrime(num) })
}

func isOdd(num int) bool { return num%2 != 0 }

func isEven(num int) bool { return num%2 == 0 }

func filterNumbersOnPredicate(numbers []int, predicate func(int) bool) []int {
	filteredNumbers := make([]int, 0, len(numbers))
	for _, num := range numbers {
		if predicate(num) {
			filteredNumbers = append(filteredNumbers, num)
		}
	}
	return filteredNumbers
}

func isPrime(number int) bool {
	if number <= 1 {
		return false
	}
	if number == 2 || number == 3 {
		return true
	}
	for i := 2; i*i <= number; i++ {
		if number%i == 0 {
			return false
		}
	}
	return true
}
