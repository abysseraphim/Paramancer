package internal

import "net/http"

type Request struct {
	Method     string
	URL        string
	Headers    map[string][]string
	Body       string
	Parameters []Parameter
}

type Parameter struct {
	Location ParamLocation
	Name     string
	Value    string
	Type     Type
}

type ParamLocation string

const (
	Query ParamLocation = "query"
	Form  ParamLocation = "form"
	JSON  ParamLocation = "json"
)

type Type string

const (
	String  Type = "string"
	Number  Type = "number"
	Boolean Type = "boolean"
	UUID    Type = "uuid"
)

type Mutation struct {
	Target   MutationTarget
	Strategy MutationStrategy
	Payload  string
}

type MutationTarget string

const (
	NameTarget  MutationTarget = "name"
	ValueTarget MutationTarget = "value"
)

type MutationStrategy string

const (
	Replace MutationStrategy = "replace"
	Append  MutationStrategy = "append"
	Prepend MutationStrategy = "prepend"
)

type Result struct {
	Request  Request
	Response *http.Response
	Error    error
}
