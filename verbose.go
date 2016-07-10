package river

import (
	"bytes"
	"fmt"
	"reflect"
	"runtime"
)

type verbose struct {
	handledPaths handledPaths
}

func (v verbose) dump() {
	var b bytes.Buffer
	fmt.Fprintln(&b)
	for _, hp := range v.handledPaths {
		fmt.Fprintf(&b, "%-4s  %-25s -> %s\n", hp.method, hp.path, nameOf(hp.handler))
	}
	logger.Println(b.String())
}

type handledPath struct {
	path    string
	method  string
	handler Handler
}

type handledPaths []handledPath

func (h *handledPaths) add(method, path string, handler Handler) {
	*h = append(*h, handledPath{path: path, method: method, handler: handler})
}

func nameOf(f interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}
