package main

import (
	"bytes"
	"text/template"
)

var apiDefinitionTmpl = template.Must(template.New("apiDefinitionTmpl").Parse(
	`type {{.RequestType}} struct {
	{{.CommonRequestFields}}
	// TODO: add {{.RequestType}} fields below
}

type {{.ResponseType}} struct {
	{{.CommonResponseFields}}
	// TODO: add {{.ResponseType}} fields below
}

func {{.MethodName}}(c *gin.Context, req *{{.RequestType}}) (resp *{{.ResponseType}}, err error) {
	{{.CommonFuncStmt}}
	// TODO: implement {{.MethodName}}
}
`))

type apiDefinitionTmplObj struct {
	RequestType          string
	CommonRequestFields  string
	ResponseType         string
	CommonResponseFields string
	MethodName           string
	CommonFuncStmt       string
}

func genHttpAPIDefinitionByTmpl(m *HttpAPIMarker, buf *bytes.Buffer, option commonHttpAPIDefinition) error {
	methodName := httpAPIMethodName(m.RequestType)
	def := apiDefinitionTmplObj{
		RequestType:          m.RequestType,
		ResponseType:         m.ResponseType,
		MethodName:           methodName,
		CommonRequestFields:  option.CommonRequestFields,
		CommonResponseFields: option.CommonResponseFields,
		CommonFuncStmt:       option.CommonFuncStmt,
	}
	return apiDefinitionTmpl.Execute(buf, &def)
}

var importTmpl = template.Must(template.New("importTmpl").Parse(
	`import {{.Alias}} "{{.PkgPath}}"
`))

type importTmplObj struct {
	PkgPath string
	Alias   string
}

func genImportByTmpl(pkgsInfo []importInfo, buf *bytes.Buffer) error {
	for _, pkg := range pkgsInfo {
		err := importTmpl.Execute(buf, importTmplObj{PkgPath: pkg.PkgPath, Alias: pkg.Alias})
		if err != nil {
			return err
		}
	}
	return nil
}
