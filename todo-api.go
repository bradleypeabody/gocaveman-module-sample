package sample

import (
	"net/http"
	"strings"

	"github.com/bradleypeabody/gocaveman"
)

// NewTodoAPIHandler makes a new TodoAPIHandler with default settings.
func NewTodoAPIHandler(todoList TodoList) *TodoAPIHandler {
	return &TodoAPIHandler{
		Prefix:   "/api/todo",
		TodoList: todoList,
	}
}

// TodoAPIHandler provides a REST API on top of TodoList
type TodoAPIHandler struct {
	Prefix   string
	TodoList TodoList
}

func (h *TodoAPIHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// see if it's our URL pattern
	if !(h.Prefix == r.URL.Path || strings.HasPrefix(r.URL.Path, h.Prefix+"/")) {
		return
	}

	err := func() error {

		var todo Todo

		// REST CRUD
		switch {

		case r.Method == "POST" && r.URL.Path == h.Prefix: // create
			err := gocaveman.JSONUnmarshalRequest(r, &todo)
			if err != nil {
				return err
			}
			newTodo, err := h.TodoList.NewTodo(todo.Text, todo.Complete)
			if err != nil {
				return err
			}
			w.Header().Set("location", h.Prefix+"/"+newTodo.TodoID)
			err = gocaveman.JSONMarshalResponse(w, 201, newTodo)
			if err != nil {
				return err
			}

		case r.Method == "GET" && strings.HasPrefix(r.URL.Path, h.Prefix+"/"): // read
			id := strings.TrimPrefix(r.URL.Path, h.Prefix+"/")
			todo, err := h.TodoList.FindTodo(id)
			if err != nil {
				return err
			}
			err = gocaveman.JSONMarshalResponse(w, 200, todo)
			if err != nil {
				return err
			}

		case (r.Method == "POST" || r.Method == "PUT" || r.Method == "PATCH") && strings.HasPrefix(r.URL.Path, h.Prefix+"/"): // update
			panic("not implemented")

		case r.Method == "DELETE" && strings.HasPrefix(r.URL.Path, h.Prefix+"/"): // delete
			panic("not implemented")

		default:
			http.NotFound(w, r)

		}

		return nil

	}()

	if err == ErrNotFound {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		gocaveman.HTTPError(w, r, err, "internal error", 500)
		return
	}

}
