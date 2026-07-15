package internal

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

var uuidRegex = regexp.MustCompile(`^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[1-5][a-fA-F0-9]{3}-[89abAB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$`)

func Parse(rawInput string, dataType InputType, scheme string) (Request, error) {
	var request Request
	var err error = nil

	if dataType == URL {
		request, err = ParseURL(rawInput)
	} else {
		request, err = ParseRawRequest(rawInput, scheme)
	}

	return request, err
}

func ParseURL(inputURL string) (Request, error) {
	var req Request
	headers := make(http.Header)

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

func ParseRawRequest(inputRequest string, scheme string) (Request, error) {
	var req Request

	reqIOreader := bufio.NewReader(strings.NewReader(inputRequest)) // new reader on raw request data
	parsedREQ, err := http.ReadRequest(reqIOreader)
	if err != nil {
		return req, err
	}

	req.URL = fmt.Sprintf("%v://%v%v", scheme, parsedREQ.Host, parsedREQ.URL.Path)
	req.Headers = parsedREQ.Header
	req.Method = parsedREQ.Method
	body, err := io.ReadAll(parsedREQ.Body)
	if err != nil {
		return req, err
	}
	defer parsedREQ.Body.Close()

	req.Body = string(body)

	req.Parameters = QueryExtractor(parsedREQ.URL.Query())

	contentType := parsedREQ.Header.Get("Content-Type")
	if strings.HasPrefix(contentType, "application/x-www-form-urlencoded") {
		extendedParams, err := FormExtractor(req.Body)
		if err == nil {
			req.Parameters = append(req.Parameters, extendedParams...)
		}
		if err != nil {
			fmt.Println("  [!]", err)
		}

	} else if strings.HasPrefix(contentType, "application/json") {
		extendedParams, err := JSONExtractor(req.Body)
		if err == nil {
			req.Parameters = append(req.Parameters, extendedParams...)
		}
		if err != nil {
			fmt.Println("  [!]", err)
		}
	}

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

func FormExtractor(body string) ([]Parameter, error) {
	var params []Parameter

	parsed, err := url.ParseQuery(body)
	if err != nil {
		return params, err
	}
	for key, vals := range parsed {
		for _, val := range vals {
			var param Parameter

			param.Location = Form
			param.Name = key
			param.Value = val
			param.Type = DetectType(val)

			params = append(params, param)
		}
	}

	return params, nil
}

func JSONExtractor(body string) ([]Parameter, error) {
	var params []Parameter
	var parsed map[string]any

	err := json.Unmarshal([]byte(body), &parsed)
	if err != nil {
		return params, err
	}

	for key, val := range parsed {
		var param Parameter

		param.Location = JSON
		param.Name = key
		param.Value = fmt.Sprint(val)
		param.Type = DetectType(param.Value)

		params = append(params, param)
	}

	return params, nil
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
