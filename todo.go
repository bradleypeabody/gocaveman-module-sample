package sample

import (
	"errors"
	"sync"

	"github.com/bradleypeabody/gouuidv6"
)

type TodoList interface {
	AllTodos() []*Todo
	ListTodos(limit int) []*Todo
	FindTodo(id string) (*Todo, error)
	NewTodo(text string, complete bool) (*Todo, error)
	DeleteTodo(id string) error
}

// NewInMemoryTodoList makes a new InMemoryTodoList with the specified configuration.
func NewInMemoryTodoList() *InMemoryTodoList {
	return &InMemoryTodoList{}
}

// Todo represents a single "to do" entry.
type Todo struct {
	TodoID   string `json:"todo_id"`
	Text     string `json:"text"`
	Complete bool   `json:"complete"`
}

var ErrNotFound = errors.New("not found")

// InMemoryTodoList is an in-memory list of "todo" items.  It manages access to them.
type InMemoryTodoList struct {
	mu       sync.RWMutex
	todoList []*Todo
}

func (tdl *InMemoryTodoList) AllTodos() []*Todo {
	return tdl.ListTodos(0)
}

func (tdl *InMemoryTodoList) ListTodos(limit int) []*Todo {
	tdl.mu.RLock()
	defer tdl.mu.RUnlock()

	if limit > 0 && len(tdl.todoList) > limit {
		return tdl.todoList[:limit]
	}

	return tdl.todoList
}

func (tdl *InMemoryTodoList) FindTodo(id string) (*Todo, error) {
	tdl.mu.RLock()
	defer tdl.mu.RUnlock()

	for _, t := range tdl.todoList {
		if t.TodoID == id {
			return t, nil
		}
	}

	return nil, ErrNotFound
}

func (tdl *InMemoryTodoList) NewTodo(text string, complete bool) (*Todo, error) {
	tdl.mu.Lock()
	defer tdl.mu.Unlock()

	t := &Todo{
		TodoID:   gouuidv6.NewB64().String(),
		Text:     text,
		Complete: complete,
	}

	tdl.todoList = append(tdl.todoList, t)

	return t, nil
}

func (tdl *InMemoryTodoList) DeleteTodo(id string) error {
	tdl.mu.Lock()
	defer tdl.mu.Unlock()

	for i, t := range tdl.todoList {
		if t.TodoID == id {

			tdl.todoList = append(tdl.todoList[:i], tdl.todoList[i+1:]...)

			return nil
		}
	}

	return ErrNotFound
}
