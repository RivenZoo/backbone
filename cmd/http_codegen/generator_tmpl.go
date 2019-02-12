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

func genHttpAPIDefinitionByTmpl(m *HttpAPIMarker, buf *bytes.Buffer, option commonHttpAPIDefinitionOption) error {
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

var apiHandlerDefineTmpl = template.Must(template.New("apiHandlerDefineTmpl").
	Delims("<?", "?>").
	Parse(`var <?.VarName?> = handler.NewRequestHandleFunc(&handler.RequestProcessor{
	NewReqFunc: func() interface{} {
		return &<?.RequestType?>{}
	},
	ProcessFunc: func(c *gin.Context, req interface{}) (resp interface{}, err error) {
		concreteReq := req.(*<?.RequestType?>)
		return <?.MethodName?>(c, concreteReq)
	},
	<? if .BodyContextKey?>BodyContextKey: <?.BodyContextKey?>, <? end ?>
	<? if .RequestDecoder?>RequestDecoder: <?.RequestDecoder?>, <? end ?>
	<? if .ResponseEncoder?>ResponseEncoder: <?.ResponseEncoder?>, <? end ?>
	<? if .ErrorEncoder?>ErrorEncoder: <?.ErrorEncoder?>, <? end ?>
	<? if .PostProcessFunc?>PostProcessFunc: <?.PostProcessFunc?>, <? end ?>
})`))

type apiHandlerDefineTmplObj struct {
	VarName     string
	RequestType string
	MethodName  string
	// optional field
	BodyContextKey  string
	RequestDecoder  string
	ResponseEncoder string
	ErrorEncoder    string
	PostProcessFunc string
}

func genAPIHandlerByTmpl(info apiHandlerDefineInfo, buf *bytes.Buffer, option commonHttpAPIHandlerOption) error {
	return apiHandlerDefineTmpl.Execute(buf, apiHandlerDefineTmplObj{
		VarName:         info.varName,
		RequestType:     info.marker.RequestType,
		MethodName:      info.apiMethodName,
		BodyContextKey:  option.BodyContextKey,
		RequestDecoder:  option.RequestDecoder,
		ResponseEncoder: option.ResponseEncoder,
		ErrorEncoder:    option.ErrorEncoder,
		PostProcessFunc: option.PostProcessFunc,
	})
}
