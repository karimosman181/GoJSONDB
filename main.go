package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"github.com/jcelliott/lumber"
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
func New(dir string, options *Options) (*Driver, error) {

	dir = filepath.Clean(dir)

	opts := Options{}

	//get options
	if options != nil {
		opts = *options
	}

	//get logger
	if opts.Logger == nil {
		opts.Logger = lumber.NewConsoleLogger((lumber.INFO))
	}

	//create driver
	driver := Driver{
		dir:     dir,
		mutexes: make(map[string]*sync.Mutex),
		log:     opts.Logger,
	}

	//check if dir exists
	if _, err := os.Stat(dir); err == nil {
		opts.Logger.Debug("Using '%s' (database already exists)\n", dir)
		return &driver, nil
	}

	//create databe if not exists
	opts.Logger.Debug("Creating the database at '%s'...\n", dir)
	return &driver, os.MkdirAll(dir, 0755)
}

// write to database
func (d *Driver) Write(collection, resource string, v interface{}) error {
	if collection == "" {
		return fmt.Errorf("Missing collection - no place to save the record!!!")
	}

	if resource == "" {
		return fmt.Errorf("Missing resource - unable to save record (no name)!!!")
	}

	mutex := d.getorCreateMutex(collection)

	//prevent changes to database until the func is complete
	mutex.lock()
	defer mutex.Unclock()

	//path to save the record
	dir := filepath.Join(d.dir, collection)
	fnlPath := filepath.Join(dir, resource+".json")
	tmpPath := fnlPath + ".tmp"

	//create dir for path
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	b, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		return err
	}

	b = append(b, byte('\n'))

	if err := ioutil.WriteFile(tmpPath, b, 0644); err != nil {
		return err
	}

	return nil
}

// read from database
func (d *Driver) Read() error {}

// read All from database
func (d *Driver) ReadAll() {}

// delete record from database
func (d *Driver) Delete() error {}

// get or create mutex if not exists
func (d *Driver) getorCreateMutex() {}

// check if dir and files exists
func stat(path string) (fi os.FileInfo, err error) {
	if fi, err = os.Stat(path); os.IsNotExist(err) {
		fi, err = os.Stat(path + ".json")
	}
	return
}

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
