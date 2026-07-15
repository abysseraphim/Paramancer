package main

import (
	"flag"
	"fmt"
	"io"
	"paramancer/internal"
	"sync"
)

const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Cyan   = "\033[36m"
	White  = "\033[37m"
	Purple = "\033[35m"
)

func main() {
	fmt.Printf(`%s

		 ██▓███   ▄▄▄       ██▀███   ▄▄▄       ███▄ ▄███▓ ▄▄▄       ███▄    █  ▄████▄  ▓█████  ██▀███  
		▓██░  ██▒▒████▄    ▓██ ▒ ██▒▒████▄    ▓██▒▀█▀ ██▒▒████▄     ██ ▀█   █ ▒██▀ ▀█  ▓█   ▀ ▓██ ▒ ██▒
		▓██░ ██▓▒▒██  ▀█▄  ▓██ ░▄█ ▒▒██  ▀█▄  ▓██    ▓██░▒██  ▀█▄  ▓██  ▀█ ██▒▒▓█    ▄ ▒███   ▓██ ░▄█ ▒
		▒██▄█▓▒ ▒░██▄▄▄▄██ ▒██▀▀█▄  ░██▄▄▄▄██ ▒██    ▒██ ░██▄▄▄▄██ ▓██▒  ▐▌██▒▒▓▓▄ ▄██▒▒▓█  ▄ ▒██▀▀█▄  
		▒██▒ ░  ░ ▓█   ▓██▒░██▓ ▒██▒ ▓█   ▓██▒▒██▒   ░██▒ ▓█   ▓██▒▒██░   ▓██░▒ ▓███▀ ░░▒████▒░██▓ ▒██▒
		▒▓▒░ ░  ░ ▒▒   ▓▒█░░ ▒▓ ░▒▓░ ▒▒   ▓▒█░░ ▒░   ░  ░ ▒▒   ▓▒█░░ ▒░   ▒ ▒ ░ ░▒ ▒  ░░░ ▒░ ░░ ▒▓ ░▒▓░
		░▒ ░       ▒   ▒▒ ░  ░▒ ░ ▒░  ▒   ▒▒ ░░  ░      ░  ▒   ▒▒ ░░ ░░   ░ ▒░  ░  ▒    ░ ░  ░  ░▒ ░ ▒░
		░░         ░   ▒     ░░   ░   ░   ▒   ░      ░     ░   ▒      ░   ░ ░ ░           ░     ░░   ░ 
		               ░  ░   ░           ░  ░       ░         ░  ░         ░ ░ ░         ░  ░   ░     
		        Parameter Mutation Engine . by abysseraphim <3                ░                     

								

%s`, Red, Reset)

	threadsFlag := flag.Int("t", 10, "Number of Workers/Threads")
	// flag.Parse()

	fmt.Printf("%s%v%s\n", Green, "  [*] starting mutation engine...", Reset)
	rawData, inputType, scheme, err := internal.ReadInput()
	if err != nil {
		fmt.Printf("%s%v%s\n\n", Red, err, Reset)
		return
	}

	request, err := internal.Parse(rawData, inputType, scheme)
	if err != nil {
		fmt.Printf("%s%v%s\n\n", Red, err, Reset)
		return
	}

	mutated, err := internal.Mutator(request)
	if err != nil {
		fmt.Printf("%s%v%s\n\n", Red, err, Reset)
		return
	}

	// for _, mReq := range mutated {
	// 	fmt.Println(mReq)
	// }

	var baseLine internal.Analysis
	orgRes, err := internal.Sender(request)
	if err != nil {
		fmt.Printf("%s%v%s\n\n", Red, "  [!] failed to fetch the original", Reset)
		return
	}

	// -------------------------------baseline
	defer orgRes.Body.Close()
	baseLine.Request = request
	baseLine.StatusCode = orgRes.StatusCode
	baseLine.Error = nil
	body, err := io.ReadAll(orgRes.Body)
	if err != nil {
		fmt.Printf("%s  [!] failed to read baseline response%s\n", Red, Reset)
	}
	baseLine.BodyLength = int64(len(body))
	baseLine.ResponseBody = string(body)
	// -------------------------------baseline

	fmt.Printf("%s%v %v% v%s\n", Green, "  [*] starting reqeuster enging with", *threadsFlag, "workers", Reset)

	jobs := make(chan internal.Request, 100)
	results := make(chan internal.Result, 100)

	var wg sync.WaitGroup

	workers := *threadsFlag
	wg.Add(workers)

	for i := 0; i < workers; i++ {
		go internal.Worker(jobs, results, &wg)
	}

	for _, mReq := range mutated {
		jobs <- mReq
	}

	close(jobs)
	go func() {
		wg.Wait()
		close(results)
	}()

	var allResults []internal.Result

	for result := range results {
		allResults = append(allResults, result)
	}
	fmt.Printf("%s  [*] collected %d responses. %s\n", Green, len(allResults), Reset)

	// fmt.Printf("%s%v%s\n", Purple, "  [*] BaseLine", Red)
	// fmt.Printf(
	// 	"%d  %d bytes  err=%v\n",
	// 	baseLine.StatusCode,
	// 	baseLine.BodyLength,
	// 	baseLine.Error,
	// )

	// fmt.Printf("%s%v%s\n", Purple, "  [*] Mutation", Red)

	analyses := internal.Analyzer(allResults)

	// for _, a := range analyses {
	// 	fmt.Printf(
	// 		"%d  %d bytes  err=%v\n",
	// 		a.StatusCode,
	// 		a.BodyLength,
	// 		a.Error,
	// 	)
	// }

	findings := internal.Detect(baseLine, analyses)

	err = internal.Reporter(request, baseLine, findings)
	if err != nil {
		fmt.Printf("%s%v%s\n\n", Red, err, Reset)
		return
	}

}
