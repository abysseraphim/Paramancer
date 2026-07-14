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

func ReadInput() (string, InputType, string, error) {

	urlFlag := flag.String("u", "", "Target URL")
	rawRequestFlag := flag.String("r", "", "Path To Your Target Request")
	schemeFlag := flag.String("s", "", "Request scheme (http or https)")

	flag.Parse()

	var dataType InputType

	if *urlFlag == "" && *rawRequestFlag == "" {
		// fmt.Println("  [!] You must specify exactly one flag with valid data")
		return "", dataType, "", errors.New("exactly one of -u or -r must be provided")

	} else if *urlFlag != "" && *rawRequestFlag != "" {
		// fmt.Println("  [!] You must specify exactly one flag with valid data")
		return "", dataType, "", errors.New("multiple input flags")
	} else if *urlFlag != "" && *schemeFlag != "" {
		return "", dataType, "", errors.New("-s can only be used with -r")
	}

	if *urlFlag != "" {
		return *urlFlag, URL, "", nil
	}

	if *rawRequestFlag != "" {

		if *schemeFlag == "http" || *schemeFlag == "https" {

			fileContent, err := os.ReadFile(*rawRequestFlag)
			if err != nil {
				return "", dataType, "", err
			}

			fileContentString := string(fileContent)
			return fileContentString, RawRequest, *schemeFlag, nil
		} else {
			return "", dataType, "", errors.New("in raw request mode, scheme is required. enter scheme with: -s http or -s https")
		}
	}

	return "", dataType, "", nil
}
