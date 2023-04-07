// Copyright 2023 Schibsted. Licensed under the terms of the Apache 2.0 license. See LICENSE in the project root.

package main

import (
	"fmt"
	"io/ioutil"

	concluder "github.com/schibsted/meeting-concluder"
)

func main() {
	content, err := ioutil.ReadFile("input.txt")
	if err != nil {
		fmt.Printf("Error reading input.txt: %v\n", err)
		return
	}

	conclusion, err := concluder.Conclude(string(content))
	if err != nil {
		fmt.Printf("Error generating conclusion: %v\n", err)
		return
	}

	conclusion += "\n"
	fmt.Print(conclusion)

	err = ioutil.WriteFile("output.txt", []byte(conclusion), 0o644)
	if err != nil {
		fmt.Printf("Error writing to output.txt: %v\n", err)
		return
	}
}
