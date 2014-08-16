package main

import "fmt"

func main() {
	for _, char := range "yeah this is good" {
		fmt.Println(char.(string))
	}
}
