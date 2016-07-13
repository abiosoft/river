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

// Dump dumps all endpoints that are being handled to the log.
func (v verbose) Dump() {
	var b bytes.Buffer
	fmt.Fprintln(&b)
	fmt.Fprintln(&b, "Endpoints")
	fmt.Fprintln(&b, "---------")
	for _, hp := range v.handledPaths {
		fmt.Fprintf(&b, "%-8s  %-25s  %s\n", hp.method, hp.path, hp.handler)
	}
	logger.Println(b.String())
}

type handledPath struct {
	path, method, handler string
}

type handledPaths []handledPath

func (h *handledPaths) add(method, path string, handler string) {
	*h = append(*h, handledPath{path: path, method: method, handler: handler})
}

func nameOf(f interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}
