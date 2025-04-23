package filter

func filterEvenNumbers(numbers []int) []int {
	return filterItemsOnAllPredicates(numbers, []Predicate[int]{isEven})
}

func filterOddNumbers(numbers []int) []int {
	return filterItemsOnAllPredicates(numbers, []Predicate[int]{isOdd})
}

func filterPrimeNumbers(numbers []int) []int {
	return filterItemsOnAllPredicates(numbers, []Predicate[int]{isPrime})
}

func filterOddPrimeNumbers(numbers []int) []int {
	return filterItemsOnAllPredicates(numbers, []Predicate[int]{isOdd, isPrime})
}

func filterEvenAndMultipleOfFiveNumbers(numbers []int) []int {
	return filterItemsOnAllPredicates(numbers, []Predicate[int]{isEven, isMultipleOfFive})
}

func filterOddAndMultipleOfThreeAndGreaterThanTen(numbers []int) []int {
	return filterItemsOnAllPredicates(numbers, []Predicate[int]{isOdd, isMultipleOfThree, isGreaterThanTen})
}

type Predicate[T any] func(T) bool

func filterItemsOnAllPredicates[T any](items []T, predicates []Predicate[T]) []T {
	filteredItems := make([]T, 0, len(items))
	for _, item := range items {
		if isAllApplicable(item, predicates) {
			filteredItems = append(filteredItems, item)
		}
	}
	return filteredItems
}

func isAllApplicable[T any](item T, predicates []Predicate[T]) bool {
	for _, predicate := range predicates {
		if !predicate(item) {
			return false
		}
	}
	return true
}

func filterItemsOnAnyPredicates[T any](items []T, predicates []Predicate[T]) []T {
	filteredItems := make([]T, 0, len(items))
	for _, item := range items {
		if isAnyApplicable(item, predicates) {
			filteredItems = append(filteredItems, item)
		}
	}
	return filteredItems
}

func isAnyApplicable[T any](item T, predicates []Predicate[T]) bool {
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
