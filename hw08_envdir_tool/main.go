package main

import (
	"log"
	"os"
)

func main() {
	args := os.Args
	if len(args) < 3 {
		log.Fatal("not enough args")
	}
	envDir, args := args[1], args[2:]
	env, err := ReadDir(envDir)
	if err != nil {
		log.Fatal(err)
	}
	oc := RunCmd(args, env)
	os.Exit(oc)
}
