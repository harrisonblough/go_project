package main

import "fmt"

func main() {
	var age [10]int
	for i := 0; i < 10 ; i++{
		age[i] = i
	}
	for i := 0; i < 10 ; i++{
		fmt.Println(age[i])
	}

}