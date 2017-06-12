package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

func check(e error) {
	if e != nil {
		fmt.Println("Error: ", e)
		panic(e)
	}
}

// Reads the contents of the file and returns
// a string containing all the contents, while
// replacing all the new line characters
func read(filename string) string {
	dat, err := ioutil.ReadFile(filename)
	check(err)
	input := string(dat)

	return strings.Replace(input, "\n", "", -1)
}
