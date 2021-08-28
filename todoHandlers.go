package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"
)

type todoHandlers struct {
	sync.Mutex
	store map[string]Todo
}

func newTodoHandlers() *todoHandlers {
	return &todoHandlers{
		store: map[string]Todo{},
	}
}

func (t *todoHandlers) postTodo(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Post Url hit")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	defer r.Body.Close()
	contentType := r.Header.Get("content-type")
	if contentType != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte(fmt.Sprintf("Invalid content-type %s, required 'application/json'", contentType)))
		return
	}
	var todo Todo
	err = json.Unmarshal(body, &todo)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	t.Lock()
	insertData(todo)
	defer t.Unlock()

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Inserted Successfully"))
}

func checkNull(str string) string {
	if str == "" {
		return "null"
	}
	return str
}

func (t *todoHandlers) getAllTodo(w http.ResponseWriter, r *http.Request) {
	var id string = checkNull(r.URL.Query().Get("id"))
	var name string = checkNull(r.URL.Query().Get("name"))
	if name != "null" {
		name = "\"" + name + "\""
	}
	var desc string = r.URL.Query().Get("description")
	if desc == "" {
		desc = "null"
	} else {
		desc = "\"%" + desc + "%\""
	}
	var priority string = checkNull(r.URL.Query().Get("priority"))
	var due string = checkNull(r.URL.Query().Get("due"))
	var complete string = checkNull(r.URL.Query().Get("completed"))
	var completion_date string = checkNull(r.URL.Query().Get("completion_date"))
	var isGetAll string = " or 1!=1"

	if id == "null" && name == "null" && desc == "null" && priority == "null" && complete == "null" && due == "null" && completion_date != "null" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid Query paramteres"))
		return
	}
	if len(r.URL.Query()) == 0 {
		isGetAll = "  or 1=1"
	}
	t.Lock()
	var due_dt string = ""
	if !(due == "null") {
		due_dt1, err := time.Parse("01-02-2006", due)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Failed to parse date.required format (MM-DD-YYYY)"))
			return
		}
		due_dt = due_dt1.Format("01-02-2006")
	}
	var comp_dt string = ""
	if !(completion_date == "null") {
		comp_dt1, err := time.Parse("01-02-2006", completion_date)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error() + "Failed to parse date.required format (MM-DD-YYYY)"))
			return
		}
		comp_dt = comp_dt1.Format("01-02-2006")
	}
	var filter_query = fmt.Sprintf("select * from todo where (%v is  null or id=%v) and "+
		"(%v is  null or name=%v) and "+
		"(%v is  null or description like %v) and "+
		"(%v is  null or priority=%v) and "+
		"(%v is  null or completed=%v) and "+
		"(%v is  null or due=str_to_date(\"%v\",\"%%m-%%d-%%Y\")) and "+
		"(%v is  null or completion_date=str_to_date(\"%v\",\"%%m-%%d-%%Y\")) %v ",
		id, id, name, name, desc, desc, priority, priority, complete, complete, due, due_dt, completion_date, comp_dt, isGetAll)

	todos := readTable(filter_query)
	t.Unlock()
	jsonbytes, err := json.Marshal(todos)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonbytes)
}

func (t *todoHandlers) mapRequest(w http.ResponseWriter, r *http.Request) {
	if !checkToken(w, r) {
		return
	}
	switch r.Method {
	case "GET":
		t.getAllTodo(w, r)
		return
	case "POST":
		fmt.Println("In Post")
		t.postTodo(w, r)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed"))
		return
	}
}

func (t *todoHandlers) mapIdByMethod(w http.ResponseWriter, r *http.Request) {
	if !checkToken(w, r) {
		return
	}

	switch r.Method {
	case "GET":
		fmt.Println("In get")
		t.getTodyById(w, r)
		return
	case "PUT":
		fmt.Println("In put")
		t.updateTodo(w, r)
		return
	case "DELETE":
		fmt.Println("In delete")
		t.deleteById(w, r)
		return
	default:
		fmt.Println("Invalid method")
		w.Write([]byte("Invalid method"))
		return
	}
}

func (t *todoHandlers) getTodyById(w http.ResponseWriter, r *http.Request) {
	vals := strings.Split(r.URL.String(), "/")
	if len(vals) != 3 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	t.Lock()
	defer t.Unlock()
	var todo Todo
	var flag bool = false
	todos := readTable("select * from todo where id=" + vals[2])
	if len(todos) > 0 {
		todo = todos[0]
		flag = true
	}

	if !flag {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	todojson, err := json.Marshal(&todo)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(todojson)
}

func (t *todoHandlers) updateTodo(w http.ResponseWriter, r *http.Request) {
	vals := strings.Split(r.URL.String(), "/")
	if len(vals) != 3 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	defer r.Body.Close()

	var todo Todo
	err = json.Unmarshal(body, &todo)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	t.Lock()
	defer t.Unlock()
	todo.Id = vals[2]
	todo = updateTodo(todo)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("updated!"))
}

func (t *todoHandlers) deleteById(w http.ResponseWriter, r *http.Request) {
	vals := strings.Split(r.URL.String(), "/")
	if len(vals) != 3 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	t.Lock()
	defer t.Unlock()

	deleteTodo(vals[2])

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("deleted!"))
}
