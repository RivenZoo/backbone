package main

import (
	"bytes"
	"go/format"
	"io"
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

type outputMerger struct {
	src   []byte
	added generatedOutputSlice
}

func (m *outputMerger) WriteTo(w io.Writer) error {
	lines := [][]byte{}
	if m.src != nil {
		lines = bytes.Split(m.src, []byte{'\n'})
	}

	buf := bytes.NewBuffer(make([]byte, 0, 1024))
	sort.Sort(m.added)

	innerIdx := 0
	for lineNo, line := range lines {
		buf.Write(line)
		buf.Write([]byte{'\n'})
		for ; innerIdx < len(m.added); innerIdx++ {
			output := m.added[innerIdx]
			if output.afterLine > lineNo+1 {
				break
			}
			buf.Write(output.buffer.Bytes())
			buf.Write([]byte{'\n'})
		}
	}
	// write remained content
	for ; innerIdx < len(m.added); innerIdx++ {
		output := m.added[innerIdx]
		buf.Write(output.buffer.Bytes())
		buf.Write([]byte{'\n'})
	}

	code, err := format.Source(buf.Bytes())
	if err != nil {
		return err
	}
	w.Write(code)
	return nil
}
