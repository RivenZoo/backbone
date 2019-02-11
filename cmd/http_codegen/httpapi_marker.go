package main

import (
	"go/ast"
	"regexp"
	"strings"
)

const httpAPIMarkerName = "HttpAPI"

var httpAPIMarkerPattern = regexp.MustCompile(`@HttpAPI\((["/\w]+)[ ]*,[ ]*([a-zA-Z][\w]*)[ ]*,[ ]*([a-zA-Z][\w]*)\)`)

type HttpAPIMarker struct {
	FileScopeMarker
	URL          string `json:"url"`
	RequestType  string `json:"request_type"`
	ResponseType string `json:"response_type"`
	commentNode  *ast.CommentGroup
}

func ParseHttpAPIMarkers(sa *SourceAst) ([]*HttpAPIMarker, error) {
	markers := make([]*HttpAPIMarker, 0)
	var err error
	for _, comments := range sa.node.Comments {
		ast.Inspect(comments, func(node ast.Node) bool {
			comment, ok := node.(*ast.Comment)
			if !ok {
				return true
			}
			var m *HttpAPIMarker
			m, err = parseHttpAPIMarkerFromCommentLine(comment.Text)
			if err != nil {
				return false
			}
			if m != nil {
				m.commentNode = comments
				markers = append(markers, m)
			}
			return true
		})
	}
	return markers, err
}

func parseHttpAPIMarkerFromCommentLine(line string) (*HttpAPIMarker, error) {
	pos := strings.Index(line, "@"+httpAPIMarkerName)
	if pos == -1 {
		return nil, nil
	}
	res := httpAPIMarkerPattern.FindAllStringSubmatch(line, -1)
	if len(res) < 1 || len(res[0]) < 4 {
		return nil, nil
	}
	matches := res[0]

	m := &HttpAPIMarker{}
	m.Name = httpAPIMarkerName
	m.URL = matches[1]
	m.RequestType = matches[2]
	m.ResponseType = matches[3]
	return m, nil
}
