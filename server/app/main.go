package main

import (
	"fmt"
	"log"
	"os"

	nc ".."
)

const (
	help          = "[USAGE]: ./main $port \ndefault port: 8989"
	incorrectArgs = "Incorrect arguments provided"
	logFail       = "Legger failed!"
)

func main() {
	file, errLog := nc.InitLogger()
	if errLog != nil {
		fmt.Println(logFail)
	}
	defer file.Close()

	// get args
	args := os.Args[1:]
	ip := "localhost"
	// port := "8989"
	port := "8181"
	if len(args) == 1 {
		port = args[0]
	} else if len(args) > 1 {
		log.Print(incorrectArgs)
		fmt.Println(help)
		os.Exit(1)
	}

	// start app listener
	if err := nc.Start(ip, port); err != nil {
		msg := fmt.Sprintf("Can't listen on port %q: %v", port, err)
		log.Print(msg)
		fmt.Fprintf(os.Stderr, msg)
		os.Exit(1)
	}
}
