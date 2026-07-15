package internal

import (
	"net/http"
	"strings"
	"sync"
	"time"
)

var client *http.Client = &http.Client{
	Timeout: 7 * time.Second,
}

func Sender(request Request) (*http.Response, error) {

	req, err := http.NewRequest(request.Method, request.URL, strings.NewReader(request.Body))
	if err != nil {
		return nil, err
	}

	req.Header = request.Headers
	if req.Header == nil {
		req.Header = make(http.Header)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/134.0.0.0 Safari/537.36 Edg/134.0.0.0")

	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func Worker(jobs <-chan Request, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()
	for req := range jobs {

		var res Result
		resp, err := Sender(req)

		res.Error = err
		res.Request = req
		res.Response = resp

		results <- res
	}
}
