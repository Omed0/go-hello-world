//package main

//import (
//	"fmt"
//	"time"
//	"unique"
//	//"crypto"
//)

//type Tasks struct {
//	id          string
//	title       string
//	description string
//	isFinished  bool
//	createdAt   time.Time
//	updatedAt   time.Time
//	deletedAt   time.Time
//}

//type user struct {
//	userId string
//	username string
//	age int
//	gender string
//	password string
//}

//func main() {
//	users :=

//	var userId string
//	var username string
//	var age int
//	var gender string
//	var password string
//	var ConfirmPassword string

//	var tasks = make([]Tasks, 0)
//	//append(tasks)

//	fmt.Println("Welcome To Task Manager, We help you to organize your tasks")
//	fmt.Print("Enter your username: ")
//	fmt.Scanln(&username)

//	fmt.Printf("Grate, Welcome %v now you can add your tasks for us and we organized for you", username)

//	fmt.Print("Now you need to register your account we ask about some information about you: ")

//	fmt.Print("Enter your age:")
//	fmt.Scanln(&age)

//	fmt.Print("Enter your gender { male | female }:")
//	fmt.Scanln(&gender)

//	fmt.Print("Enter your Password:")
//	fmt.Scanln(&password)

//	fmt.Print("Enter your Confirm Password:")
//	fmt.Scanln(&ConfirmPassword)

//	if password != ConfirmPassword {
//		fmt.Println("Your Passwords does not same match make again")
//		//continue
//	}
//	tempId := unique.Make(username)
//	userId = tempId.Value()

//}

package main

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

type User struct {
	UserId   string
	Username string
	Email    string
	Age      int
}

func GetFakeUsers() []User {

	users := []User{
		{
			UserId:   "1001",
			Username: "alice_wonder",
			Email:    "test@gmail.com",
			Age:      25,
		}, {
			UserId:   "1002",
			Username: "bob_builder",
			Email:    "bob_builder@gmail.com",
			Age:      32,
		}, {
			UserId:   "1003",
			Username: "charlie_fox",
			Email:    "charlie_fox@gmail.com",
			Age:      28,
		}, {
			UserId:   "1004",
			Username: "dana_smith",
			Email:    "dana_smith@gmail.com",
			Age:      22,
		}, {
			UserId:   "1005",
			Username: "eve_black",
			Email:    "eve_black@gmail.com",
			Age:      29,
		},
	}
	return users
}

type Key struct {
	Path, Country string
}

func main() {
	users := GetFakeUsers()
	fmt.Println("Users data:")

	var wg sync.WaitGroup
	for _, user := range users {
		wg.Add(1) // Increment the WaitGroup counter for each user
		go func(email, username string) {
			defer wg.Done()
			sendEmail(email, fmt.Sprintf("Hello %s, welcome to our service!", username))
		}(user.Email, user.Username)
	}

	jsonData, err := json.Marshal(users)
	if err != nil {
		return
	} else {
		fmt.Println(string(jsonData))
		fmt.Println("Users data has been successfully marshaled to JSON format.")
	}

	wg.Wait() // Wait for all goroutines to complete
}

func sendEmail(email, message string) {
	time.Sleep(2 * time.Second) // Simulate a delay for sending email
	// Simulate sending an email
	fmt.Printf("Sending email to %s with message: %s\n", email, message)
	// In a real application, you would use an SMTP client or similar to send the email
}
