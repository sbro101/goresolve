package main

// this is to test if the module is working

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/sbro101/goresolve"
)

func main() {
	url := os.Args[1]

	rd := goresolve.Hostname(url, "1.1.1.1")

	json, err := json.MarshalIndent(rd, "", "  ")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%s\n", json)
}
