package main

import (
	"fmt"
	"paramancer/internal"
)

func main() {
	fmt.Println("Hello World!")
	rawData, inputType, scheme, err := internal.ReadInput()
	if err != nil {
		fmt.Println(err)
		return
	}

	request, err := internal.Parse(rawData, inputType, scheme)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(request)

}
