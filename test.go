package main

import "fmt"

func main() {
	// Reverse a
	a := []string{"1", "2"}
	for i, j := 0, len(a) - 1; i < j; i, j = i+1, j-1 {
		a[i], a[j] = a[j], a[i]
	}

	fmt.Println(a)
}
