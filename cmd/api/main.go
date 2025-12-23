package main

import "github.com/Ablebil/lathi-be/cmd/bootstrap"

func main() {
	if err := bootstrap.Start(); err != nil {
		panic(err)
	}
}
