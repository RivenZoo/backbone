package main

import (
	"bytes"
	"text/template"
)

var apiDefinitionTmpl = template.Must(template.New("apiDefinitionTmpl").Parse(
	`type {{.RequestType}} struct {
	// TODO: add {{.RequestType}} fields below
}

type {{.ResponseType}} struct {
	// TODO: add {{.ResponseType}} fields below
}

func {{.MethodName}}(c *gin.Context, req *{{.RequestType}}) (resp *{{.ResponseType}}, err error) {
	// TODO: implement {{.MethodName}}
}
`))

type apiDefinitionTmplObj struct {
	RequestType  string
	ResponseType string
	MethodName   string
}

func genHttpAPIDefinitionByTmpl(m *HttpAPIMarker, buf *bytes.Buffer) error {
	methodName := httpAPIMethodName(m.RequestType)
	def := apiDefinitionTmplObj{
		RequestType:  m.RequestType,
		ResponseType: m.ResponseType,
		MethodName:   methodName,
	}
	return apiDefinitionTmpl.Execute(buf, &def)
}

var importTmpl = template.Must(template.New("importTmpl").Parse(
	`import "{{.PkgPath}}"
`))

type importTmplObj struct {
	PkgPath string
}

func genImportByTmpl(pkgs []string, buf *bytes.Buffer) error {
	for _, pkg := range pkgs {
		err := importTmpl.Execute(buf, importTmplObj{PkgPath: pkg})
		if err != nil {
			return err
		}
	}
	return nil
}
