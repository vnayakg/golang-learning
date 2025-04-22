package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

type Vertex struct {
	x, y int
}

func (v Vertex) abs() float64 {
	return math.Sqrt(float64(v.x*v.x + v.y*v.y))
}

func abs(v Vertex) float64 {
	return math.Sqrt(float64(v.x*v.x + v.y*v.y))
}

// scale with pointer receiver
func (v *Vertex) scale(factor int) {
	v.x *= factor
	v.y *= factor
}

// scale with value receiver
func (v Vertex) scaleWithValue(factor int) Vertex {
	v.scale(factor)
	return v
}

// scale with just function
func scale(v *Vertex, factor int) {
	v.scale(factor)
}

// method on non-struct type
type Float float64

func (f Float) abs() float64 {
	if f >= 0 {
		return float64(f)
	}
	return float64(-f)
}

type Abser interface {
	abs() float64
}

type IPAddr [4]byte

func (ip IPAddr) String() string {
	s := make([]string, len(ip))
	for i, val := range ip {
		s[i] = strconv.Itoa(int(val))
	}
	return strings.Join(s, ".")
}

// Errors
type MyError struct {
	message string
}

func (e MyError) Error() string {
	return fmt.Sprintf("message: %v", e.message)
}

func run() error {
	return MyError{"something bad happened"}
}

type ErrNegativeSqrt float64

func (e ErrNegativeSqrt) Error() string {
	return fmt.Sprintf("cannot Sqrt negative number: %v", float64(e))
}

func Sqrt(x float64) (float64, error) {
	if x < 0 {
		return 0, ErrNegativeSqrt(x)
	}
	return math.Sqrt(x), nil
}

func main() {
	v := Vertex{1, 1}
	fmt.Println(v.abs())
	fmt.Println(abs(v))

	fmt.Println(Float(1).abs())

	v.scale(10)
	fmt.Println(v)
	fmt.Println(v.scaleWithValue(10))

	scale(&v, 99)
	fmt.Println(v)

	//interfaces
	var a Abser
	a = &v
	fmt.Println(a.abs())
	a = Float(1)
	fmt.Println(a.abs())

	t, ok := a.(Vertex) //type assertion
	fmt.Println(t, ok)

	switch s := a.(type) {
	case Vertex:
		fmt.Println("type of a is Vertex")
	case Float:
		fmt.Println("type of a is Float")
	default:
		fmt.Printf("%v %T\n", s, s)
	}

	//printing ip addresses exercise
	hosts := map[string]IPAddr{
		"loopback": {127, 0, 0, 1},
		"someDNS":  {1, 2, 3, 4},
	}
	for name, ip := range hosts {
		fmt.Printf("%v: %v\n", name, ip)
	}

	//Errors
	res := run()

	if res != nil {
		fmt.Println(res)
	}

	// sqrt error exercise
	fmt.Println(Sqrt(2))
	fmt.Println(Sqrt(-2))
}
