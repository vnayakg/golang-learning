package main

import "fmt"

type Vertex struct {
	x, y int
}

func main() {
	// pointers
	i := 2
	p := &i
	fmt.Println("p value: ", p)
	fmt.Println("p pointing to: ", *p)
	fmt.Printf("type of p is %T\n", p)

	//struct
	v := Vertex{1, 2}
	fmt.Println(v)
	v.x = 100
	fmt.Println(v.x)

	//pointer to struct
	pointerToV := &v
	pointerToV.y = 200
	fmt.Println(v.x, v.y)
	anotherPointerToStruct := &Vertex{-1, -1}
	fmt.Println(*anotherPointerToStruct, anotherPointerToStruct.x, anotherPointerToStruct.y)
	anotherV := Vertex{}
	fmt.Println(anotherV)

	//arrays
	array := [5]int{1, 2, 3, 4, 5}
	fmt.Println(array)
	slicedArray := array[1:3]
	fmt.Println(slicedArray)
	array[3] = 1
	fmt.Printf("%v len=%d cap=%d\n", slicedArray, len(slicedArray), cap(slicedArray))
	colors := []string{"red", "green"}
	fmt.Println(colors)
}
