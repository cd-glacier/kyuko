package main

import (
	"fmt"

	kyuko "github.com/g-hyoga/kyuko/go"
)

func main() {

	err := kyuko.Exec()
	if err != nil {
		fmt.Printf("error!!\n%s", err)
	}

}
