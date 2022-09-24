package main

import (
	awssecrets "awssecrets/secrets"
	"encoding/json"
	"fmt"
	"log"
)

func main() {
	result, err := awssecrets.GetValues()
	if err != nil {
		log.Fatal(err)
	}
	jsonString, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(jsonString))
}
