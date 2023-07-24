package main

import (
	"encoding/json"
	"fmt"
)

type Address struct {
	City    string      `json:"city"`
	State   string      `json:"state"`
	Country string      `json:"country"`
	Pincode json.Number `json:"pincode"`
}

type User struct {
	Name    string      `json:"name"`
	Age     json.Number `json:"age"`
	Contact string      `json:"contact"`
	Company string      `json:"company"`
	Address Address     `json:"address"`
}

func main() {
	dir := "./"

	db, err := New(dir, nil)
	if err != nil {
		fmt.Println("Error", err)
	}
}
