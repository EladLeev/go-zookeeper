//
// A simple exaple on how to use the package.
// Get, Children, Exist,
//
package main

import (
	"fmt"
	"time"

	"github.com/samuel/go-zookeeper/zk"
)

const (
	zkPath = "/test"
	zkNode = "localhost"
)

func getZnode(c *zk.Conn, path string) string {
	// Use `Get` to simply get data from ZK.
	// `Get` gets a string path, and return the data as []byte
	data, _, err := c.Get(path)
	if err != nil {
		panic(err)
	}
	// Transform the []byte into string
	s := string(data[:])
	return s
}

func getChildren(c *zk.Conn, path string) []string {
	// Use `Children` to get a slice of strings with all the children of provided path.
	data, _, err := c.Children(zkPath)
	if err != nil {
		panic(err)
	}
	return data
}

func checkZnode(c *zk.Conn, path string) error {
	// Use `Exists` to check if zNode exist.
	// Retrun value is bool
	_, _, err := c.Exists(zkPath)
	return err
}

func main() {
	// Connect to ZK, print if there is an error, and close the connection at the end
	c, _, err := zk.Connect([]string{zkNode}, time.Second)
	defer c.Close()
	if err != nil {
		panic(err)
	}
	// Example #1
	// Just a simple Get
	// Get the data of spesific zNode
	myData := getZnode(c, zkPath)
	fmt.Printf("The data is: %v\n", myData)

	// Example #2
	// Get zNode Children and for each one of them,
	// check if the znode exist and then print the key and the data.
	myChld := getChildren(c, zkPath)
	fmt.Printf("The key are: %v\n", myChld)
	for _, key := range myChld {
		chldPath := zkPath + "/" + key
		err := checkZnode(c, chldPath)
		if err != nil {
			panic(err)
		}
		fmt.Printf("The key is: %v\n", key)
		fmt.Printf("The Data is: %v\n", getZnode(c, chldPath))
	}

	// Example #3
	// Set watch on znode, and wait for event.
	// Catch the event and use switch - case to do something, or just print the event name.
	data, _, ch, err := c.GetW(zkPath)
	if err != nil {
		panic(err)
	}
	fmt.Printf("The current data is: %v.\nI'm watching.\n", string(data[:]))
	for event := range ch {
		switch event.Type {
		case 1:
			fmt.Printf("Node created.")
			// Do somthing
		case 2:
			fmt.Printf("Node deleted!")
			// Do somthing
		case 3:
			fmt.Printf("Node changed.")
			// Do somthing
		case 4:
			fmt.Printf("Node children changed.")
			// Do somthing
		default:
			fmt.Printf("Something happen.")
		}
		// Or just ...
		fmt.Println(event.Type)

	}
}
