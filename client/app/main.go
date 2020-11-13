package main

import (
	"fmt"
	"os"

	client ".."
)

func main() {
	ip := "localhost"
	port := "8181"
	o, err := client.StartClient(ip, port)
	if err != nil {
		fmt.Printf("Cannot connect to %v:%v\n", ip, port)
		os.Exit(1)
	}
	// client.Test(o)
	// o := &client.IO{}
	if err := client.Main(o); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer o.Conn.Close()
}
