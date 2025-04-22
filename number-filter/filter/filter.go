package filter

func filterEvenNumbers(numbers []int) []int {
	return filterItemsOnAllPredicates(numbers, []func(int) bool{isEven})
}

func filterOddNumbers(numbers []int) []int {
	return filterItemsOnAllPredicates(numbers, []func(int) bool{isOdd})
}

func filterPrimeNumbers(numbers []int) []int {
	return filterItemsOnAllPredicates(numbers, []func(int) bool{isPrime})
}

func filterOddPrimeNumbers(numbers []int) []int {
	return filterItemsOnAllPredicates(numbers, []func(int) bool{isOdd, isPrime})
}

func filterEvenAndMultipleOfFiveNumbers(numbers []int) []int {
	return filterItemsOnAllPredicates(numbers, []func(int) bool{isEven, isMultipleOfFive})
}

func filterOddAndMultipleOfThreeAndGreaterThanTen(numbers []int) []int {
	return filterItemsOnAllPredicates(numbers, []func(int) bool{isOdd, isMultipleOfThree, isGreaterThanTen})
}

func filterItemsOnAllPredicates[T any](items []T, predicates []func(T) bool) []T {
	filteredItems := make([]T, 0, len(items))
	for _, item := range items {
		if isAllApplicable(item, predicates) {
			filteredItems = append(filteredItems, item)
		}
	}
	return filteredItems
}

func isAllApplicable[T any](item T, predicates []func(T) bool) bool {
	for _, predicate := range predicates {
		if !predicate(item) {
			return false
		}
	}
	return true
}

func filterItemsOnAnyPredicates[T any](items []T, predicates []func(T) bool) []T {
	filteredItems := make([]T, 0, len(items))
	for _, item := range items {
		if isAnyApplicable(item, predicates) {
			filteredItems = append(filteredItems, item)
		}
	}
	return filteredItems
}

func isAnyApplicable[T any](item T, predicates []func(T) bool) bool {
	for _, predicate := range predicates {
		if predicate(item) {
			return true
		}
	}
	return false
}

func isOdd(num int) bool { return num%2 != 0 }

func isEven(num int) bool { return num%2 == 0 }

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
