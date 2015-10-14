package swugger

type ServiceDoc struct {
	Doc string
	Consumes []string
	Produces []string
}

type MethodDoc struct {
	Operation string
	Doc string
	Notes string
	Params []ParamDoc
	Reads interface{}
	Writes interface{}
}

type ParamDoc struct {
	Type string
	Name string
	Doc string
	DataType string
}
