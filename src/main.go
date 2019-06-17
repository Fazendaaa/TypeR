package main

import (
	"fmt"
	"os"
	"os/user"

	"./repl"
)

func main() {
	user, err := user.Current()

	if nil != err {
		panic(err)
	}

	fmt.Printf("Hello %s! This is TypeR programming language!\n", user.Username)
	fmt.Printf("Fell free to type in commands\n")

	repl.Start(os.Stdin, os.Stdout)
}
