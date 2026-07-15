package internal

import (
	"encoding/json"
	"net/url"
	"strconv"
)

func Builder(original Request, target Parameter, mutation Mutation) (Request, error) {
	newRequest := original

	switch target.Location {
	case Query:
		err := BuildQuery(&newRequest, target, mutation)
		if err != nil {
			return newRequest, err
		}
	case Form:
		err := BuildForm(&newRequest, target, mutation)
		if err != nil {
			return newRequest, err
		}
	case JSON:
		err := BuildJSON(&newRequest, target, mutation)
		if err != nil {
			return newRequest, err
		}
	}

	return newRequest, nil
}

func BuildQuery(req *Request, target Parameter, mutation Mutation) error {
	queryParams := make(url.Values)
	for _, param := range req.Parameters {
		if param.Location != Query {
			continue
		}

		if param == target {
			// do mutation here
			if mutation.Target == NameTarget {
				// muatue name
				if mutation.Strategy == Append {
					// append to name
					param.Name = param.Name + mutation.Payload
				} else if mutation.Strategy == Replace {
					// replace name
					param.Name = mutation.Payload
				} else if mutation.Strategy == Prepend {
					// prepend to name
					param.Name = mutation.Payload + param.Name
				}

			} else if mutation.Target == ValueTarget {
				// mutate value
				if mutation.Strategy == Append {
					// append to value
					param.Value = param.Value + mutation.Payload
				} else if mutation.Strategy == Replace {
					// replace value
					param.Value = mutation.Payload
				} else if mutation.Strategy == Prepend {
					// prepend to value
					param.Value = mutation.Payload + param.Value
				}
			}
		}

		queryParams.Add(param.Name, param.Value)
	}

	if queryParams.Encode() != "" {
		req.URL = req.URL + "?" + queryParams.Encode()
	}

	return nil
}

func BuildForm(req *Request, target Parameter, mutation Mutation) error {
	queryParams := make(url.Values)

	for _, param := range req.Parameters {
		if param.Location != Form {
			continue
		}

		if param == target {
			if mutation.Target == NameTarget {
				if mutation.Strategy == Append {
					param.Name = param.Name + mutation.Payload
				} else if mutation.Strategy == Replace {
					param.Name = mutation.Payload
				} else if mutation.Strategy == Prepend {
					param.Name = mutation.Payload + param.Name
				}

			} else if mutation.Target == ValueTarget {
				if mutation.Strategy == Append {
					param.Value = param.Value + mutation.Payload
				} else if mutation.Strategy == Replace {
					param.Value = mutation.Payload
				} else if mutation.Strategy == Prepend {
					param.Value = mutation.Payload + param.Value
				}
			}
		}

		queryParams.Add(param.Name, param.Value)
	}

	if queryParams.Encode() != "" {
		req.Body = queryParams.Encode()
	}

	return nil
}

func BuildJSON(req *Request, target Parameter, mutation Mutation) error {
	jsonParams := make(map[string]any)

	for _, param := range req.Parameters {
		if param.Location != JSON {
			continue
		}

		if param == target {
			if mutation.Target == NameTarget {
				if mutation.Strategy == Append {
					param.Name = param.Name + mutation.Payload
				} else if mutation.Strategy == Replace {
					param.Name = mutation.Payload
				} else if mutation.Strategy == Prepend {
					param.Name = mutation.Payload + param.Name
				}
			} else if mutation.Target == ValueTarget {
				if mutation.Strategy == Append {
					param.Value = param.Value + mutation.Payload
				} else if mutation.Strategy == Replace {
					param.Value = mutation.Payload
				} else if mutation.Strategy == Prepend {
					param.Value = mutation.Payload + param.Value

				}
			}
		}

		var value any
		switch param.Type {
		case Number:
			num, err := strconv.Atoi(param.Value)
			if err != nil {
				value = param.Value
			}
			if err == nil {
				value = num
			}

		case Boolean:
			booli, err := strconv.ParseBool(param.Value)
			if err != nil {
				value = param.Value
			}
			if err == nil {
				value = booli
			}

		case UUID:
			value = param.Value
		default:
			value = param.Value
		}

		jsonParams[param.Name] = value
	}

	jsonBody, err := json.Marshal(jsonParams)
	if err != nil {
		return err
	}

	req.Body = string(jsonBody)
	return nil
}
