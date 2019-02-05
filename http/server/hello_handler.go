package server

import (
	"fmt"
	"net/http"
)

type helloHandler struct {
	welcome string
}

func (h helloHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(fmt.Sprintf("<html><body>%s</body></html>", h.welcome)))
}
