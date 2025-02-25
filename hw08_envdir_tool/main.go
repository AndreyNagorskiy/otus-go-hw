package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args

	if len(args) < 3 {
		panic("wrong command arguments count")
	}

	env, err := ReadDir(args[1])
	if err != nil {
		fmt.Println(err)
	}

	RunCmd(args[2:], env)
}
