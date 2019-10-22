package main

import (
	"fmt"
	"os"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("usage: prof {ref} {command}")
		os.Exit(1)
	}

	ref := os.Args[1]
	command := os.Args[2:]
	fmt.Printf("%v %v\n", ref, command)

}
