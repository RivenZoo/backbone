package main

import (
	"bytes"
	"fmt"
	"github.com/RivenZoo/backbone/logger"
	"io"
	"io/ioutil"
	"strings"
)

const registerRouterFuncName = "Post"

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

type commonInitRouterStmtOption struct {
	MiddlewareNames []string
}

type httpAPIGeneratorOption struct {
	apiDefineFileImports  []importInfo
	apiHandlerFileImports []importInfo
	initRouterImports     []importInfo
	commonAPIDefinition   commonHttpAPIDefinitionOption
	commonHttpAPIHandler  commonHttpAPIHandlerOption
	commonInitRouter      commonInitRouterStmtOption
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

	routerInitFile        string
	routerInitFileContent []byte
	routerInitSource      *SourceAst

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

func (g *HttpAPIGenerator) parseAPIHandlerFile() error {
	data, err := readSource(g.handlerDefineFile)
	if err != nil {
		return err
	}
	if data == nil {
		// no file content
		return nil
	}

	g.handlerFileContent = data

	sa, err := ParseSourceCode(g.handlerDefineFile, bytes.NewReader(data))
	if err != nil {
		logger.Errorf("ParseSourceCode %s error %v", g.handlerDefineFile, err)
		return err
	}
	g.handlerSource = sa
	return nil
}

func (g *HttpAPIGenerator) filterUndeclaredAPIHandler() []apiHandlerDefineInfo {
	if g.handlerSource == nil {
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
	if err := g.parseAPIHandlerFile(); err != nil {
		return
	}
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
		end := 2
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

func httpRouterInitFilename(pkgName string) string {
	return fmt.Sprintf("%s_urls.go", pkgName)
}

func httpRouterInitFuncName() string {
	return "registerUrls"
}

func (g *HttpAPIGenerator) GenInitHttpAPIRouter() {
	pkgName := g.packageName()
	g.routerInitFile = httpRouterInitFilename(pkgName)

	if g.parseInitRouterFile() != nil {
		return
	}
	stmtInfos := g.filterUnRegisterAPI(httpRouterInitFuncName())
	if len(stmtInfos) == 0 {
		return
	}

	g.addInitRouterImports()
	unImported := g.option.initRouterImports
	if g.routerInitSource != nil {
		unImported = filterUnImportedPackage(g.routerInitSource.node.Imports, g.option.initRouterImports)
	}
	g.genInitHttpAPIRouterOutput(stmtInfos, unImported)
}

func (g *HttpAPIGenerator) OutputInitHttpAPIRouter(w io.Writer) {
	merger := outputMerger{
		src:   g.routerInitFileContent, //maybe nil
		added: g.routerInitOutput,
	}
	merger.WriteTo(w)
}

func (g *HttpAPIGenerator) parseInitRouterFile() error {
	data, err := readSource(g.routerInitFile)
	if err != nil {
		return err
	}
	if data == nil {
		// no file content
		return nil
	}

	g.routerInitFileContent = data

	sa, err := ParseSourceCode(g.routerInitFile, bytes.NewReader(data))
	if err != nil {
		logger.Errorf("ParseSourceCode %s error %v", g.routerInitFile, err)
		return err
	}
	g.routerInitSource = sa
	return nil
}

type initRouterStmtInfo struct {
	marker          *HttpAPIMarker
	handlerFuncName string
	afterLine       int
}

func (g *HttpAPIGenerator) filterUnRegisterAPI(initFunc string) []initRouterStmtInfo {
	if g.routerInitSource == nil {
		ret := make([]initRouterStmtInfo, 0, len(g.markers))
		for _, m := range g.markers {
			ret = append(ret, initRouterStmtInfo{
				marker:          m,
				afterLine:       5, // package,space,import,space,func
				handlerFuncName: httpAPIMethodName(m.RequestType),
			})
		}
		return ret
	}

	return filterUnInitRouters(g.routerInitSource, initFunc, registerRouterFuncName, g.markers)
}

func (g *HttpAPIGenerator) addInitRouterImports() {
	requiredImports := []importInfo{
		{"github.com/gin-gonic/gin", ""},
	}
	g.option.initRouterImports = mergeImports(g.option.initRouterImports, requiredImports)
}

func (g *HttpAPIGenerator) genInitHttpAPIRouterOutput(stmtInfos []initRouterStmtInfo, unImported []importInfo) {
	if len(g.routerInitFileContent) <= 0 {
		// empty file, add package declare
		g.routerInitOutput = append(g.routerInitOutput, generatedOutput{
			buffer:    bytes.NewBufferString(fmt.Sprintf("package %s\n", g.source.node.Name.Name)),
			afterLine: 1,
		})
	}
	if len(unImported) > 0 {
		end := 2
		if g.routerInitSource != nil {
			// after package stmt
			end = g.routerInitSource.fSet.Position(g.routerInitSource.node.Name.End()).Line
		}
		output := generatedOutput{
			buffer:    genImports(unImported),
			afterLine: end,
		}
		g.routerInitOutput = append(g.routerInitOutput, output)
	}

	initFuncName := httpRouterInitFuncName()
	genInitFunc := true
	if g.routerInitSource != nil {
		fnDecl := filterGlobalFunc(g.routerInitSource.node.Decls, initFuncName)
		if fnDecl != nil {
			genInitFunc = false
		}
	}

	initRouterVarName := "engine"
	if genInitFunc {
		// begin new init func define
		buf := bytes.NewBuffer(make([]byte, 0))
		closeFunc := genFuncDefine(initFuncName, []string{"engine *gin.Engine"}, buf)
		end := 3 // after package, import
		if len(g.routerInitOutput) > 0 {
			end = g.routerInitOutput[len(g.routerInitOutput)-1].afterLine + 1
		}
		g.routerInitOutput = append(g.routerInitOutput, generatedOutput{
			buffer:    buf,
			afterLine: end,
		})
		defer func() {
			// end new init func define
			buf := bytes.NewBuffer(make([]byte, 0))
			closeFunc(buf)
			g.routerInitOutput = append(g.routerInitOutput, generatedOutput{
				buffer:    buf,
				afterLine: g.routerInitOutput[len(g.routerInitOutput)-1].afterLine + 1,
			})
		}()
	}
	initRouterStmtFn := func(stmtInfo initRouterStmtInfo) *bytes.Buffer {
		buf := bytes.NewBuffer(make([]byte, 0))
		genInitRouterStmtByTmpl(stmtInfo, initRouterVarName, registerRouterFuncName, buf, g.option.commonInitRouter)
		return buf
	}
	for _, info := range stmtInfos {
		buf := initRouterStmtFn(info)
		output := generatedOutput{
			afterLine: info.afterLine,
			buffer:    buf,
		}
		g.routerInitOutput = append(g.routerInitOutput, output)
	}
}

func (g *HttpAPIGenerator) packageName() string {
	return g.source.node.Name.Name
}

func mergeImports(srcImports, addImports []importInfo) []importInfo {
	findExistImport := func(pkg string) bool {
		for _, imp := range srcImports {
			if imp.PkgPath == pkg {
				return true
			}
		}
		return false
	}
	imports := srcImports
	for _, imp := range addImports {
		if !findExistImport(imp.PkgPath) {
			imports = append(imports, imp)
		}
	}
	return imports
}
