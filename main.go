package main

import "fmt"

func main() {
	var hello string = "Hello, World!"
	nameProgram := "Go"

	fmt.Printf("Welcome to the %v program!\n", nameProgram) // This program demonstrates basic variable declaration and function usage in Go.
	fmt.Println(hello)

	var x uint = 10
	fmt.Println("Value of x:", x)

	var y int = 10
	fmt.Println("Value of y:", y)

	fmt.Println("Value of total is:", sum(x, y))
}

func sum(a uint, b int) uint {
	return a + uint(b) // The function sum takes a uint and an int, adds them, and returns a uint.
}
