package filter

import (
	"reflect"
	"testing"
)

func TestItShouldFilterEvenNumbers(t *testing.T) {
	testCases := []struct {
		input    []int
		expected []int
	}{
		{[]int{}, []int{}},
		{[]int{0, 2, 4, 6}, []int{0, 2, 4, 6}},
		{[]int{1, 3, 5}, []int{}},
		{[]int{1, 2, 3, 4}, []int{2, 4}},
	}

	for _, testCase := range testCases {
		actual := filterEvenNumbers(testCase.input)

		if !reflect.DeepEqual(actual, testCase.expected) {

			t.Errorf("For input %v, expected %v, but got %v", testCase.input, testCase.expected, actual)
		}
	}
}

func TestItShouldFilterOddNumbers(t *testing.T) {
	testCases := []struct {
		input    []int
		expected []int
	}{
		{[]int{}, []int{}},
		{[]int{1, 3, 5}, []int{1, 3, 5}},
		{[]int{0, 2, 4}, []int{}},
		{[]int{1, 2, 3, 4}, []int{1, 3}},
	}

	for _, testCase := range testCases {
		actual := filterOddNumbers(testCase.input)

		if !reflect.DeepEqual(actual, testCase.expected) {

			t.Errorf("For input %v, expected %v, but got %v", testCase.input, testCase.expected, actual)
		}
	}
}

func TestItShouldFilterPrimeNumbers(t *testing.T) {
	testCases := []struct {
		input    []int
		expected []int
	}{
		{[]int{}, []int{}},
		{[]int{-1, 1, 4}, []int{}},
		{[]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 37, 101}, []int{2, 3, 5, 7, 37, 101}},
	}

	for _, testCase := range testCases {
		actual := filterPrimeNumbers(testCase.input)

		if !reflect.DeepEqual(actual, testCase.expected) {

			t.Errorf("For input %v, expected %v, but got %v", testCase.input, testCase.expected, actual)
		}
	}
}

func TestItShouldFilterOddPrimeNumbers(t *testing.T) {
	testCases := []struct {
		input    []int
		expected []int
	}{
		{[]int{}, []int{}},
		{[]int{-1, 1, 4}, []int{}},
		{[]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 37, 101}, []int{3, 5, 7, 37, 101}},
	}

	for _, testCase := range testCases {
		actual := filterOddPrimeNumbers(testCase.input)

		if !reflect.DeepEqual(actual, testCase.expected) {

			t.Errorf("For input %v, expected %v, but got %v", testCase.input, testCase.expected, actual)
		}
	}
}

func TestItShouldFilterEvenMultipleOfFiveNumbers(t *testing.T) {
	testCases := []struct {
		input    []int
		expected []int
	}{
		{[]int{}, []int{}},
		{[]int{1, 5, 15, 4}, []int{}},
		{[]int{10, 20, 30, 5, 8}, []int{10, 20, 30}},
	}

	for _, testCase := range testCases {
		actual := filterEvenAndMultipleOfFiveNumbers(testCase.input)

		if !reflect.DeepEqual(actual, testCase.expected) {

			t.Errorf("For input %v, expected %v, but got %v", testCase.input, testCase.expected, actual)
		}
	}
}
