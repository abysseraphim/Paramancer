package internal

import (
	"errors"
	"flag"
	"os"
)

type InputType string

const (
	URL        InputType = "url"
	RawRequest InputType = "raw"
)

func ReadInput() (string, InputType, error) {

	urlFlag := flag.String("u", "", "Target URL")
	rawRequestFlag := flag.String("r", "", "Path To Your Target Request")

	flag.Parse()

	var dataType InputType

	if *urlFlag == "" && *rawRequestFlag == "" {
		// fmt.Println("  [!] You must specify exactly one flag with valid data")
		return "", dataType, errors.New("empty or invalid flag data")

	} else if *urlFlag != "" && *rawRequestFlag != "" {
		// fmt.Println("  [!] You must specify exactly one flag with valid data")
		return "", dataType, errors.New("multiple input flags")
	}

	if *urlFlag != "" {
		return *urlFlag, URL, nil
	}

	if *rawRequestFlag != "" {
		fileContent, err := os.ReadFile(*rawRequestFlag)
		if err != nil {
			return "", dataType, err
		}

		fileContentString := string(fileContent)
		return fileContentString, RawRequest, nil
	}

	return "", dataType, nil
}
