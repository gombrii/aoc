package internal

import "fmt"

// NOTE: Generates project scaffolding, which is really just the shared folder with some nice to haves:
// - input
// - exit
// - render
// - go.mod
// Maybe the init AOC command should also take a single string with the mod name, like go mod init
func GenAoc(module string ) error {
	//Generate shared folder with all it contains and go.mod file, I guess with the go version the user has
	fmt.Println("NOT IMPLEMENTED:", module)
	return nil
}
