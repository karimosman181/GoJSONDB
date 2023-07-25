package main

import (
	"encoding/json"
	"fmt"
	"sync"
)

const Version = "1.0.1"

type (
	Logger interface {
		Fatal(string, ...interface{})
		Error(string, ...interface{})
		Warn(string, ...interface{})
		Info(string, ...interface{})
		Debug(string, ...interface{})
		Trace(string, ...interface{})
	}

	Driver struct {
		mutex   sync.Mutex
		mutexes map[string]*sync.Mutex
		dir     string
		log     Logger
	}
)

type Options struct {
	Logger
}

type Address struct {
	City    string      `json:"city"`
	State   string      `json:"state"`
	Country string      `json:"country"`
	Zipcode json.Number `json:"zipcode"`
}

type User struct {
	Name    string      `json:"name"`
	Age     json.Number `json:"age"`
	Contact string      `json:"contact"`
	Company string      `json:"company"`
	Address Address     `json:"address"`
}

// intialize and create new database
func New() {}

// write to database
func Write() error {}

// read from database
func Read() error {}

// read All from database
func ReadAll() {}

// delete record from database
func Delete() error {}

// get or create mutex if not exists
func getorCreateMutex() {}

func main() {
	dir := "./"

	db, err := New(dir, nil)
	if err != nil {
		fmt.Println("Error", err)
	}

	employees := []User{
		{"John", "23", "John@email.j", "JOhn tech", Address{"John city", "John state", "John Country", "0001"}},
		{"Alex", "33", "Alex@email.a", "Alex tech", Address{"Alex city", "Alex state", "Alex Country", "0002"}},
		{"Max", "33", "Max@email.m", "Max tech", Address{"Max city", "Max state", "Max Country", "0003"}},
	}

	//loop over employees and insert to db
	for _, value := range employees {
		//insert employee to db
		db.Write("users", value.Name, User{
			Name:    value.Name,
			Age:     value.Age,
			Contact: value.Contact,
			Company: value.Company,
			Address: value.Address,
		})
	}

	//get all users records
	records, err := db.ReadAll("users")
	if err != nil {
		fmt.Println("Error", err)
	}
	fmt.Println(records)

	allusers := []User{}

	for _, f := range records {

		employeeFound := User{}
		if err := json.Unmarshal([]byte(f), &employeeFound); err != nil {
			fmt.Println("Error", err)
		}

		allusers = append(allusers, employeeFound)

	}

	fmt.Println((allusers))

	//delete user
	// if err := db.Delete("user","john"); err != nil {
	// 	fmt.Println("Error", err)
	// }

	// if err := db.Delete("user",""); err != nil {
	// 	fmt.Println("Error", err)
	// }

}
