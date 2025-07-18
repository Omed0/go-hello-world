package main

import (
	"fmt"
	"io"
	"os"
	"runtime"
)

func main() {
	var hello string = "Hello, World!"
	nameProgram := "Go"

	fmt.Printf("Welcome to the %v program!\n", nameProgram) // This program demonstrates basic variable declaration and function usage in Go.
	fmt.Printf("Message: %s\n", hello)                      // This line prints the message stored in the variable hello.
	//fmt.Println(hello)

	var x uint = 10
	fmt.Println("Value of x:", x)

	var y int = 10
	fmt.Println("Value of y:", y)

	fmt.Println("Value of total is:", sum(x, y))
	defineOS() // Call the function to define the operating system.

	readFile("example.txt") // Call the function to read a file named "example.txt".

	pointer(10)

	fmt.Println("End of the program")
}

func sum(a uint, b int) uint {
	return a + uint(b) // The function sum takes a uint and an int, adds them, and returns a uint.
}

func defineOS() {
	os := runtime.GOOS // This function retrieves the current operating system.
	switch os {
	case "windows":
		fmt.Println("Running on Windows")
	case "linux":
		fmt.Println("Running on Linux")
	case "darwin":
		fmt.Println("Running on macOS")
	default:
		fmt.Printf("Running on %s\n", os)
	}
	runtime.GC()
	// This line forces garbage collection, which is not typically necessary in Go but can be used for demonstration purposes.
}

func readFile(filename string) {
	f, err := os.Open(filename)

	if err != nil {
		//log.Fatal(err)
		fmt.Println(err)
	}
	defer f.Close()

	//read all of file
	t, err := io.ReadAll(f)

	if err != nil {
		//log.Fatal(err)
		fmt.Println(err)
	}

	fmt.Printf("File content: %s\n", string(t)) // This line prints the content of the file.
	if len(t) == 0 {
		fmt.Println("File is empty")
	} else {
		fmt.Println("File is not empty")
	}
	fmt.Println("File read successfully")
}

func pointer(value int) {
	x := value
	y := &x // This line creates a pointer to the variable x.

	fmt.Println("Value of x:", x)
	fmt.Println("Pointer to x:", y) // This line prints the memory address of x.

	x += 10
	*y += 10
	fmt.Println(x, *y) // this print value x and value pointer, and because we work with reference address must both return same value and +20

	var username string
	fmt.Scan(username) // this just print empty because it is empty
	fmt.Print("Put username: ")
	fmt.Scanln(&username) // but this let u input value to reference memory of username or pointer

	fmt.Printf("Welcome Again %v\n", username)
}
