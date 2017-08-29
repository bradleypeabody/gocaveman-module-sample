package sample

import (
	"net/http"
	"path"

	"github.com/bradleypeabody/gocaveman"
)

func NewTodoIncludeFS() http.FileSystem {
	return gocaveman.NewHTTPFuncFS(func(name string) (http.File, error) {
		name = path.Clean("/" + name)
		if name == "/todo/mini-list.gohtml" {
			return gocaveman.NewHTTPBytesFile(name, modTime, todoMiniListIncludeBytes), nil
		}
		return nil, ErrNotFound
	})
}

var todoMiniListIncludeBytes = []byte(`
{{range}}
put our mini list thing here
{{end}}
`)
