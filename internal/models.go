package internal

import "net/http"

type Request struct {
	Method     string
	URL        string
	Headers    http.Header
	Body       string
	Parameters []Parameter

	Mutation Mutation
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
	Parameter string
	Location  ParamLocation
	Target    MutationTarget
	Strategy  MutationStrategy
	Payload   string
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

type Analysis struct {
	Request      Request
	StatusCode   int
	BodyLength   int64
	Error        error
	ResponseBody string
}

type Finding struct {
	Request  Request
	Reason   string
	Severity Severity
	Evidence string
	Mutation Mutation
}

type Severity string

const (
	Info     Severity = "info"
	Low      Severity = "low"
	Medium   Severity = "medium"
	High     Severity = "high"
	Critical Severity = "critical"
	Unknown  Severity = "unknown"
)

type Report struct {
	Target string

	Baseline Analysis

	Findings []FinalResult
}

type FinalResult struct {
	Severity Severity
	Reason   string
	Evidence string

	Parameter         string
	ParamLocation     ParamLocation
	ParamTargetedPart MutationTarget
	MutationStrategy  MutationStrategy
	Payload           string
}
