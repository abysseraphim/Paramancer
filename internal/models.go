package internal

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
