package internal

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

var uuidRegex = regexp.MustCompile(`^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[1-5][a-fA-F0-9]{3}-[89abAB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$`)

func Parse(rawInput string, dataType InputType) (Request, error) {
	var request Request
	var err error = nil

	if dataType == URL {
		request, err = ParseURL(rawInput)
	} else {
		request, err = ParseRawRequest(rawInput)
	}

	return request, err
}

func ParseURL(inputURL string) (Request, error) {
	var req Request
	headers := make(map[string][]string)
	// we gonna need some headers when we send the request, headers like user agent.

	parsedURL, err := url.Parse(inputURL)
	if err != nil {
		return req, err
	}

	req.Method = "GET"
	req.Headers = headers
	req.Body = ""
	req.URL = fmt.Sprintf("%v://%v%v", parsedURL.Scheme, parsedURL.Host, parsedURL.Path)
	req.Parameters = QueryExtractor(parsedURL.Query())

	return req, nil
}

func ParseRawRequest(inputRequest string) (Request, error) {
	var req Request

	return req, nil
}

func QueryExtractor(values url.Values) []Parameter {
	var params []Parameter

	for key, vals := range values {
		for _, val := range vals {
			var param Parameter

			param.Location = Query
			param.Name = key
			param.Value = val
			param.Type = DetectType(val)

			params = append(params, param)
		}
	}

	return params
}

func DetectType(value string) Type {
	if uuidRegex.MatchString(value) {
		return UUID
	} else if strings.ToLower(value) == "true" || strings.ToLower(value) == "false" {
		return Boolean
	} else if _, err := strconv.Atoi(value); err == nil {
		return Number
	}

	return String
}
