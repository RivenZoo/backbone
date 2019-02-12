package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"sort"
)

type generatedOutput struct {
	buffer    *bytes.Buffer
	afterLine int
}

type generatedOutputSlice []generatedOutput

func (o generatedOutputSlice) Len() int           { return len(o) }
func (o generatedOutputSlice) Swap(i, j int)      { o[i], o[j] = o[j], o[i] }
func (o generatedOutputSlice) Less(i, j int) bool { return o[i].afterLine < o[j].afterLine }

type importInfo struct {
	PkgPath string
	Alias   string
}

type commonHttpAPIDefinition struct {
	CommonRequestFields  string
	CommonResponseFields string
	CommonFuncStmt       string
}

type httpAPIGeneratorOption struct {
	imports             []importInfo
	commonAPIDefinition commonHttpAPIDefinition
}

type httpAPIGenerator struct {
	option            httpAPIGeneratorOption
	srcFile           string
	srcContent        []byte
	source            *SourceAst
	markers           []*HttpAPIMarker
	handlerDefineFile string
	routerInitFile    string

	// generated output
	sourceFileOutput []generatedOutput // request/response type declare, handle func declare
	handlerOutput    []generatedOutput // gin handle func definition
	routerInitOutput []generatedOutput // register gin handler
}

func newHttpAPIGenerator(option httpAPIGeneratorOption) *httpAPIGenerator {
	ret := &httpAPIGenerator{
		option:  option,
		markers: []*HttpAPIMarker{},
	}
	return ret
}

func (g *httpAPIGenerator) parseFile(srcFile string) error {
	sa, err := ParseSourceFile(srcFile)
	if err != nil {
		return err
	}
	g.source = sa
	g.srcFile = srcFile
	return nil
}

func (g *httpAPIGenerator) parseCode(fname string, code io.Reader) error {
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

func (g *httpAPIGenerator) parseHttpAPIMarkers() error {
	m, err := ParseHttpAPIMarkers(g.source)
	if err != nil {
		return err
	}
	g.markers = m
	return nil
}

func (g *httpAPIGenerator) genHttpAPIDeclare() {
	unImported := filterUnImportedPackage(g.source.node.Imports, g.option.imports)
	declaredFuncs := filterDelcaredFuncNames(g.source.node.Decls)
	unDeclaredMarkers := filterFuncUndeclaredHttpAPIMarkers(g.markers, declaredFuncs)
	g.genSourceOutput(unDeclaredMarkers, unImported)
}

func (g *httpAPIGenerator) outputAPIDeclare(w io.Writer) error {
	if g.srcContent == nil {
		data, err := ioutil.ReadFile(g.srcFile)
		if err != nil {
			return err
		}
		g.srcContent = data
	}
	sort.Sort(generatedOutputSlice(g.sourceFileOutput))
	lines := bytes.Split(g.srcContent, []byte{'\n'})
	innerIdx := 0
	for lineNo, line := range lines {
		w.Write(line)
		w.Write([]byte{'\n'})
		for ; innerIdx < len(g.sourceFileOutput); innerIdx++ {
			output := g.sourceFileOutput[innerIdx]
			if output.afterLine > lineNo+1 {
				break
			}
			w.Write(output.buffer.Bytes())
			w.Write([]byte{'\n'})
		}
	}
	return nil
}

func (g *httpAPIGenerator) genSourceOutput(unDeclaredMarkers []*HttpAPIMarker, unImported []importInfo) {
	g.sourceFileOutput = make([]generatedOutput, 0, len(unDeclaredMarkers))
	importDecl := func(pkgs []importInfo) *bytes.Buffer {
		buf := bytes.NewBuffer(make([]byte, 0))
		genImportByTmpl(pkgs, buf)
		return buf
	}
	if len(unImported) > 0 {
		end := g.source.fSet.Position(g.source.node.Name.End()).Line
		output := generatedOutput{
			buffer:    importDecl(unImported),
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

func (g *httpAPIGenerator) genHttpAPIHandler() {

}

func (g *httpAPIGenerator) genInitHttpAPIRouter() {

}
