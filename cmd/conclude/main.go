package main

import (
	"fmt"
	"io/ioutil"
	"strings"

	concluder "github.schibsted.io/alexander-fet-rodseth/hackday-meeting-concluder"
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

	if !strings.HasSuffix(conclusion, "\n") {
		conclusion += "\n"
	}

	fmt.Print(conclusion)

	err = ioutil.WriteFile("output.txt", []byte(conclusion), 0o644)
	if err != nil {
		fmt.Printf("Error writing to output.txt: %v\n", err)
		return
	}
}
