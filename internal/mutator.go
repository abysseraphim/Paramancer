package internal

import "fmt"

func Mutator(req Request) ([]Request, error) {
	var requests []Request

	targets := []MutationTarget{NameTarget, ValueTarget}
	strategies := []MutationStrategy{Append, Replace, Prepend}
	payloads := []string{"'", "\"", "[", "]", "%00", "%09", "%0D", "%0A", "\t", "\r", "\n", "\\", "{{7*7}}", "${7*7}", "[]", "null", "NaN"}

	for _, param := range req.Parameters {
		for _, target := range targets {
			for _, strategy := range strategies {
				for _, payload := range payloads {
					var mutation Mutation

					mutation.Target = target
					mutation.Strategy = strategy
					mutation.Payload = payload
					mutation.Parameter = param.Name
					mutation.Location = param.Location

					forged, err := Builder(req, param, mutation)
					if err != nil {
						fmt.Printf("\033[31m  [!] %v\033[0m\n", err)
						continue
					}

					requests = append(requests, forged)
				}
			}
		}
	}

	return requests, nil
}
