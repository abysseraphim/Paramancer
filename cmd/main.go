package main

import (
	"fmt"
	"paramancer/internal"
)

func main() {
	fmt.Println("Hello World!")
	rawData, inputType, err := internal.ReadInput()
	if err != nil {
		fmt.Println(err)
		return
	}

	request, err := internal.Parse(rawData, inputType)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(request)
}
