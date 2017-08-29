package sample

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path"
	"testing"

	"github.com/bradleypeabody/gocaveman"
	"github.com/stretchr/testify/assert"
)

// func TestTodoServer(t *testing.T) {

// 	t.Logf("TestTodoServer")

// 	// skip if flag not passed - print log statement with suggested command line to run including timeout

// 	// setup server

// 	// run the server, print URL to test

// }

func TestTodo(t *testing.T) {

	assert := assert.New(t)
	// setup server

	// do some API tests against it with HTTP client

	// FIXME: we should provide a blank empty page and stuff like this as part of the main Caveman package - modules will be doing this all over the place as part of thier test
	incFS := gocaveman.NewHTTPFuncFS(func(name string) (http.File, error) {
		name = path.Clean("/" + name)
		if name == "/main-page.gohtml" {
			return gocaveman.NewHTTPBytesFile(name, modTime, []byte(`<!DOCTYPE html>
<html>
<head>
	<title></title>
</head>
<body>

{{block "body" .}}{{end}}

</body>
</html>`)), nil
		}
		return nil, ErrNotFound
	})

	todoList := NewInMemoryTodoList()

	// templates
	rendererHandler := gocaveman.NewDefaultRenderer(
		NewTodoViewFS(),
		gocaveman.NewHTTPStackedFileSystem(incFS, NewTodoIncludeFS()),
	)

	h := gocaveman.BuildHandlerChain(gocaveman.HandlerList{
		gocaveman.NewGzipHandler(), // FIXME: this really should not be called GzipHandler - the entire chain and list mechanism depends on it
		NewTodoAPIHandler(todoList),
		rendererHandler,
		http.NotFoundHandler(),
	})

	s := httptest.NewServer(h)
	defer s.Close()

	client := s.Client()

	req, err := http.NewRequest("POST", s.URL+"/api/todo", bytes.NewBufferString(`{"text":"My TODO item here"}`))
	assert.Nil(err)
	req.Header.Set("content-type", "application/json")
	res, err := client.Do(req)
	assert.Nil(err)
	defer res.Body.Close()
	assert.Equal(201, res.StatusCode)

	loc := res.Header.Get("location")

	t.Logf("loc: %q", loc)

	url := s.URL + loc
	t.Logf("getting URL: %q", url)
	req, err = http.NewRequest("GET", url, nil)
	assert.Nil(err)
	req.Header.Set("content-type", "application/json")
	res, err = client.Do(req)
	assert.Nil(err)
	defer res.Body.Close()
	assert.Equal(200, res.StatusCode)

	b, err := ioutil.ReadAll(res.Body)
	assert.Nil(err)
	t.Logf("Response body: %q", b)

	bs := string(b)

	assert.Contains(bs, `"text":"My TODO item here"`)

}

// func setupTodoServer() {

// 	// FIXME: moving to docs/developer-overview.gohtml
// 	// integration points:
// 	// - initialization
// 	// - request handler
// 	// - view fs layer
// 	// - include fs layer
// 	// - request context

// 	gzipHandler := gocaveman.NewGzipHandler()

// 	// static stuff we need assigned to each context
// 	ctxMapHandler := gocaveman.NewCtxMapHandler(map[string]interface{}{
// 	// "menus": menus,
// 	})

// 	// templates
// 	viewsFs := gocaveman.NewHTTPStackedFileSystem(
// 		afero.NewHttpFs(afero.NewBasePathFs(rootFs, viper.GetString("views-path"))),
// 		editorFs,
// 	)
// 	rendererHandler := gocaveman.NewDefaultRenderer(
// 		viewsFs,
// 		afero.NewHttpFs(afero.NewBasePathFs(rootFs, viper.GetString("includes-path"))),
// 	)

// 	// static files
// 	staticHandler := gocaveman.NewStaticFileServer(afero.NewBasePathFs(rootFs, viper.GetString("static-path")), "/")

// 	// assemble our handlers with the appropriate sequence/hierarchy
// 	// FIXME: we can make something that checks for SetNextHandler and calls it automatically, so we can just keep a slice with the sequence
// 	h := gzipHandler.SetNextHandler(
// 		ctxMapHandler.SetNextHandler(
// 			gocaveman.NewHandlerChain(
// 				defaultRedirectHandler,
// 				rendererHandler,
// 				staticHandler,
// 				http.NotFoundHandler(),
// 			)))

// 	log.Fatal(httpServer.ListenAndServe())

// }
