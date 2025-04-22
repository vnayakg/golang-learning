package filter

func filterEvenNumbers(numbers []int) []int {
	return filterNumbersOnAllPredicates(numbers, []func(int) bool{isEven})
}

func filterOddNumbers(numbers []int) []int {
	return filterNumbersOnAllPredicates(numbers, []func(int) bool{isOdd})
}

func filterPrimeNumbers(numbers []int) []int {
	return filterNumbersOnAllPredicates(numbers, []func(int) bool{isPrime})
}

func filterOddPrimeNumbers(numbers []int) []int {
	return filterNumbersOnAllPredicates(numbers, []func(int) bool{isOdd, isPrime})
}

func filterEvenAndMultipleOfFiveNumbers(numbers []int) []int {
	return filterNumbersOnAllPredicates(numbers, []func(int) bool{isEven, isMultipleOfFive})
}

func filterOddAndMultipleOfThreeAndGreaterThanTen(numbers []int) []int {
	return filterNumbersOnAllPredicates(numbers, []func(int) bool{isOdd, isMultipleOfThree, isGreaterThanTen})
}

func isOdd(num int) bool { return num%2 != 0 }

func isEven(num int) bool { return num%2 == 0 }

func filterNumbersOnAllPredicates(numbers []int, predicates []func(int) bool) []int {
	filteredNumbers := make([]int, 0, len(numbers))
	for _, num := range numbers {
		if isAllApplicable(num, predicates) {
			filteredNumbers = append(filteredNumbers, num)
		}
	}
	return filteredNumbers
}

func isAllApplicable(number int, predicates []func(int) bool) bool {
	for _, predicate := range predicates {
		if !predicate(number) {
			return false
		}
	}
	return true
}

func filterNumbersOnAnyPredicates(numbers []int, predicates []func(int) bool) []int {
	filteredNumbers := make([]int, 0, len(numbers))
	for _, num := range numbers {
		if isAnyApplicable(num, predicates) {
			filteredNumbers = append(filteredNumbers, num)
		}
	}
	return filteredNumbers
}

func isAnyApplicable(number int, predicates []func(int) bool) bool {
	for _, predicate := range predicates {
		if predicate(number) {
			return true
		}
	}
	return false
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

func isMultipleOfFive(number int) bool {
	return number%5 == 0
}

func isMultipleOfThree(number int) bool {
	return number%3 == 0
}

func isGreaterThanTen(number int) bool {
	return number > 10
}
