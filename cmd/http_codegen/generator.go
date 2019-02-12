package main

import (
	"bytes"
	"fmt"
	"github.com/RivenZoo/backbone/logger"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

type importInfo struct {
	PkgPath string
	Alias   string
}

type commonHttpAPIDefinitionOption struct {
	CommonRequestFields  string
	CommonResponseFields string
	CommonFuncStmt       string
}

type commonHttpAPIHandlerOption struct {
	BodyContextKey  string
	RequestDecoder  string
	ResponseEncoder string
	ErrorEncoder    string
	PostProcessFunc string
}

type httpAPIGeneratorOption struct {
	apiDefineFileImports  []importInfo
	apiHandlerFileImports []importInfo
	initRouterImports     []importInfo
	commonAPIDefinition   commonHttpAPIDefinitionOption
	commonHttpAPIHandler  commonHttpAPIHandlerOption
}

type HttpAPIGenerator struct {
	option     httpAPIGeneratorOption
	srcFile    string
	srcContent []byte
	source     *SourceAst
	markers    []*HttpAPIMarker

	handlerDefineFile  string
	handlerSource      *SourceAst
	handlerFileContent []byte

	routerInitFile string

	// generated output
	sourceFileOutput []generatedOutput // request/response type declare, handle func declare
	handlerOutput    []generatedOutput // gin handle func definition
	routerInitOutput []generatedOutput // register gin handler
}

func newHttpAPIGenerator(option httpAPIGeneratorOption) *HttpAPIGenerator {
	ret := &HttpAPIGenerator{
		option:  option,
		markers: []*HttpAPIMarker{},
	}
	return ret
}

func (g *HttpAPIGenerator) ParseFile(srcFile string) error {
	sa, err := ParseSourceFile(srcFile)
	if err != nil {
		return err
	}
	g.source = sa
	g.srcFile = srcFile
	return nil
}

func (g *HttpAPIGenerator) ParseCode(fname string, code io.Reader) error {
	data, err := ioutil.ReadAll(code)
	if err != nil {
		return err
	}
	g.srcContent = data

	sa, err := ParseSourceCode(fname, bytes.NewReader(data))
	if err != nil {
		return err
	}
	g.source = sa
	g.srcFile = fname
	return nil
}

func (g *HttpAPIGenerator) ParseHttpAPIMarkers() error {
	m, err := ParseHttpAPIMarkers(g.source)
	if err != nil {
		return err
	}
	g.markers = m
	return nil
}

func (g *HttpAPIGenerator) GenHttpAPIDeclare() {
	unImported := filterUnImportedPackage(g.source.node.Imports, g.option.apiDefineFileImports)
	declaredFuncs := filterDelcaredFuncNames(g.source.node.Decls)
	unDeclaredMarkers := filterFuncUndeclaredHttpAPIMarkers(g.markers, declaredFuncs)
	g.genAPIDeclareOutput(unDeclaredMarkers, unImported)
}

func (g *HttpAPIGenerator) OutputAPIDeclare(w io.Writer) error {
	if g.srcContent == nil {
		data, err := ioutil.ReadFile(g.srcFile)
		if err != nil {
			return err
		}
		g.srcContent = data
	}
	merger := outputMerger{
		src:   g.srcContent,
		added: g.sourceFileOutput,
	}
	return merger.WriteTo(w)
}

func genImports(pkgs []importInfo) *bytes.Buffer {
	buf := bytes.NewBuffer(make([]byte, 0))
	genImportByTmpl(pkgs, buf)
	return buf
}

func httpAPIMethodName(requestTypeName string) string {
	return fmt.Sprintf("handle%s", strings.Title(requestTypeName))
}

func (g *HttpAPIGenerator) genAPIDeclareOutput(unDeclaredMarkers []*HttpAPIMarker, unImported []importInfo) {
	g.sourceFileOutput = make([]generatedOutput, 0, len(unDeclaredMarkers))
	if len(unImported) > 0 {
		end := g.source.fSet.Position(g.source.node.Name.End()).Line
		output := generatedOutput{
			buffer:    genImports(unImported),
			afterLine: end,
		}
		g.sourceFileOutput = append(g.sourceFileOutput, output)
	}

	apiTypeDeclare := func(m *HttpAPIMarker) *bytes.Buffer {
		buf := bytes.NewBuffer(make([]byte, 0))
		genHttpAPIDefinitionByTmpl(m, buf, g.option.commonAPIDefinition)
		return buf
	}
	for _, m := range unDeclaredMarkers {
		output := generatedOutput{
			buffer:    apiTypeDeclare(m),
			afterLine: m.EndLine,
		}
		g.sourceFileOutput = append(g.sourceFileOutput, output)
	}
}

type apiHandlerDefineInfo struct {
	marker        *HttpAPIMarker
	afterLine     int
	varName       string
	apiMethodName string
}

func httpAPIHandlerVarName(reqType string) string {
	return fmt.Sprintf("%sHandler", reqType)
}

func (g *HttpAPIGenerator) filterUndeclaredAPIHandler() []apiHandlerDefineInfo {
	data, err := ioutil.ReadFile(g.handlerDefineFile)
	if err != nil && os.IsNotExist(err) {
		ret := []apiHandlerDefineInfo{}
		for i, m := range g.markers {
			ret = append(ret, apiHandlerDefineInfo{
				marker:        m,
				afterLine:     i + 1,
				varName:       httpAPIHandlerVarName(m.RequestType),
				apiMethodName: httpAPIMethodName(m.RequestType),
			})
		}
		return ret
	}
	if err != nil {
		logger.Errorf("open file %s error %v", g.handlerDefineFile, err)
		return nil
	}
	g.handlerFileContent = data

	sa, err := ParseSourceCode(g.handlerDefineFile, bytes.NewReader(data))
	if err != nil {
		logger.Errorf("ParseSourceCode %s error %v", g.handlerDefineFile, err)
		return nil
	}
	g.handlerSource = sa

	return filterUndeclaredHandlers(g.handlerSource, g.markers)
}

func apiHandlerFileName(srcFilename string) string {
	return fmt.Sprintf("%s_handlers.go", srcFilename)
}

func (g *HttpAPIGenerator) addAPIHandlerImports() {
	requiredImports := []importInfo{
		{"github.com/RivenZoo/backbone/http/handler", ""},
		{"github.com/gin-gonic/gin", ""},
	}

	g.option.apiHandlerFileImports = mergeImports(g.option.apiHandlerFileImports, requiredImports)
}

func (g *HttpAPIGenerator) GenHttpAPIHandler() {
	g.handlerDefineFile = apiHandlerFileName(g.source.filePath)
	handlerDefineInfos := g.filterUndeclaredAPIHandler()
	if len(handlerDefineInfos) == 0 {
		return
	}

	g.addAPIHandlerImports()
	unImported := g.option.apiHandlerFileImports
	if g.handlerSource != nil {
		unImported = filterUnImportedPackage(g.handlerSource.node.Imports, g.option.apiHandlerFileImports)
	}
	g.genHttpAPIHandlerOutput(handlerDefineInfos, unImported)
}

func (g *HttpAPIGenerator) OutputAPIHandler(w io.Writer) error {
	merger := outputMerger{
		src:   g.handlerFileContent, // maybe it's nil
		added: g.handlerOutput,
	}
	return merger.WriteTo(w)
}

func (g *HttpAPIGenerator) genHttpAPIHandlerOutput(handlerDefineInfos []apiHandlerDefineInfo,
	unImported []importInfo) {
	if len(g.handlerFileContent) <= 0 {
		// empty file, add package declare
		g.handlerOutput = append(g.handlerOutput, generatedOutput{
			buffer:    bytes.NewBufferString(fmt.Sprintf("package %s\n", g.source.node.Name.Name)),
			afterLine: 1,
		})
	}
	if len(unImported) > 0 {
		end := 1
		if g.handlerSource != nil {
			end = g.handlerSource.fSet.Position(g.handlerSource.node.Name.End()).Line
		}
		output := generatedOutput{
			buffer:    genImports(unImported),
			afterLine: end,
		}
		g.handlerOutput = append(g.handlerOutput, output)
	}
	apiHandlerDefine := func(info apiHandlerDefineInfo) *bytes.Buffer {
		buf := bytes.NewBuffer(make([]byte, 0))
		genAPIHandlerByTmpl(info, buf, g.option.commonHttpAPIHandler)
		return buf
	}
	for _, info := range handlerDefineInfos {
		buf := apiHandlerDefine(info)
		output := generatedOutput{
			afterLine: info.afterLine,
			buffer:    buf,
		}
		g.handlerOutput = append(g.handlerOutput, output)
	}
}

func (g *httpAPIGenerator) genInitHttpAPIRouter() {

}
