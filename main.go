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
)

// "encoding/json"
// "fmt"

type User struct {
	UserId   string
	Username string
	Age      int
}

func GetFakeUsers() []User {

	users := []User{
		{
			UserId:   "1001",
			Username: "alice_wonder",
			Age:      25,
		}, {
			UserId:   "1002",
			Username: "bob_builder",
			Age:      32,
		}, {
			UserId:   "1003",
			Username: "charlie_fox",
			Age:      28,
		}, {
			UserId:   "1004",
			Username: "dana_smith",
			Age:      22,
		}, {
			UserId:   "1005",
			Username: "eve_black",
			Age:      29,
		},
	}
	return users
}

type Key struct {
	Path, Country string
}

func main() {
	hits := make(map[Key]int)

	k := Key{"main.go", "US"}
	s := Key{"main.go", "US"}
	hits[k]++
	hits[s]++

	//get value from key main.g
	for key, value := range hits {
		fmt.Println("Key:", key.Path, "Value:", value)
	}

	users := GetFakeUsers()

	jsonData, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		return
	} else {
		fmt.Println(string(jsonData))
		fmt.Println("Users data has been successfully marshaled to JSON format.")
	}

}
