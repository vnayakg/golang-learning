package main

import (
	"fmt"
	"math"
)

const (
	Big   = 1 << 100
	Small = Big >> 99
	Pi    = 3.14
)

func main() {
	fmt.Printf("hello")
	printZeroValue()
	typeConversion()
	printConstant()
	fmt.Println(Small)
	fmt.Println(needInt(Small))
	fmt.Println(needFloat(Small))
	fmt.Println(needFloat(Big))
	fmt.Println(getSumAndAbsDiff(5, 6))
}

func getSumAndAbsDiff(x, y int) (int, int) {
	return x + y, int(math.Abs(float64(x - y)))
}

func printZeroValue() {
	var i int
	var f float64
	var c complex64
	var b bool
	var s string
	fmt.Printf("%v %v %v %q %v\n", i, f, b, s, c)
}

func typeConversion() {
	var x, y int = 3, 5
	var f float64 = math.Sqrt(float64(x*x + y*y))
	var z uint = uint(f)
	fmt.Println(x, y, z, f)
}

func printConstant() {
	const World = "世界"
	fmt.Println("Hello", World)
	fmt.Println("Happy", Pi, "Day")

	const Truth = true
	fmt.Println("Go rules?", Truth)
}

func needInt(x int) int { return x*10 + 1 }

func needFloat(x float64) float64 {
	return x * 0.1
}
