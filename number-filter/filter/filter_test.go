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
