package sample

import (
	"net/http"
	"path"
	"time"

	"github.com/bradleypeabody/gocaveman"
)

func NewTodoViewFS() http.FileSystem {
	return gocaveman.NewHTTPFuncFS(func(name string) (http.File, error) {
		name = path.Clean("/" + name)
		if name == "/todo.gohtml" {
			return gocaveman.NewHTTPBytesFile(name, modTime, todoListViewBytes), nil
		}
		if name == "/todo/edit.gohtml" {
			return gocaveman.NewHTTPBytesFile(name, modTime, todoEditViewBytes), nil
		}
		return nil, ErrNotFound
	})
}

var modTime = time.Now()

var todoListViewBytes = []byte(`{{template "main-page.gohtml" .}}
{{define "body"}}
List view
{{end}}
`)

var todoEditViewBytes = []byte(`{{template "main-page.gohtml" .}}
{{define "body"}}
Edit view
{{end}}
`)
