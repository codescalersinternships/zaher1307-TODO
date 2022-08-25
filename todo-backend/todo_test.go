package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var s = Server{}

func initDatabase(str string) {

	var err error

	if s.db, err = gorm.Open(sqlite.Open(str+DB_FILE), &gorm.Config{}); err != nil {
		fmt.Println("Error")
		return
	}
	s.db.AutoMigrate(&Todo{})

}

func removeDatabase(str string) {
	os.RemoveAll(str + DB_FILE)
}

func TestGetTodoItemsList(t *testing.T) {

	t.Run("try to retrieve all todo items", func(t *testing.T) {

		initDatabase("1")
		defer removeDatabase("1")

		request := httptest.NewRequest(http.MethodGet, "localhost:8080/todolist", nil)
		response := httptest.NewRecorder()

		want := []Todo{}
		s.db.Find(&want)

		s.GetTodoItemsList(response, request)

		got := []Todo{}
		json.NewDecoder(response.Body).Decode(&got)

		if !reflect.DeepEqual(want, got) {
			t.Errorf("\ngot:\n%v\nwant:\n%v\n", got, want)
		}

	})

}

func TestGetTodoItem(t *testing.T) {

	t.Run("try to retrieve a specific todo item", func(t *testing.T) {

		initDatabase("2")
		defer removeDatabase("2")

		request := httptest.NewRequest(http.MethodGet, "localhost:8080/todolist/1", nil)
		response := httptest.NewRecorder()

		want := Todo{ID: 1, TodoItem: "study", Completed: false}
		s.db.Create(&want)

		s.GetTodoItem(response, request)

		var got Todo
		json.NewDecoder(response.Body).Decode(&got)

		if !reflect.DeepEqual(want, got) {
			t.Errorf("\ngot:\n%v\nwant:\n%v\n", got, want)
		}

	})

}

func TestAddTodoItem(t *testing.T) {

	t.Run("try to get a bad request", func(t *testing.T) {

		initDatabase("3")
		defer removeDatabase("3")

		postRequestBody := strings.NewReader(`"ihjkldshl,asdfklskl[sdasdafjssdji()]"false}`)
		postRequest := httptest.NewRequest(http.MethodPost, "localhost:8080/todolist", postRequestBody)
		postResponse := httptest.NewRecorder()

		s.AddTodoItem(postResponse, postRequest)
		got := postResponse.Result().Status
		want := "400 Bad Request"

		if !reflect.DeepEqual(want, got) {
			t.Errorf("\ngot:\n%v\nwant:\n%v\n", got, want)
		}

	})

	t.Run("try to get a bad request again", func(t *testing.T) {

		initDatabase("3")
		defer removeDatabase("3")

		postRequestBody := strings.NewReader(`{"invalid":"invalid", "invalid":"study", "complete":false}`)
		postRequest := httptest.NewRequest(http.MethodPost, "localhost:8080/todolist", postRequestBody)
		postResponse := httptest.NewRecorder()

		s.AddTodoItem(postResponse, postRequest)
		got := postResponse.Result().Status
		want := "400 Bad Request"

		if !reflect.DeepEqual(want, got) {
			t.Errorf("\ngot:\n%v\nwant:\n%v\n", got, want)
		}

	})

}

func TestUpdataTodoItem(t *testing.T) {

	t.Run("try to get a bad request in modification", func(t *testing.T) {

		initDatabase("3")
		defer removeDatabase("3")

		patchRequestBody := strings.NewReader(`"ihjkldshl,asdfklskl[sdasdafjssdji()]"false}`)
		patchRequest := httptest.NewRequest(http.MethodPatch, "localhost:8080/todolist", patchRequestBody)
		patchResponse := httptest.NewRecorder()

		s.UpdateTodoItem(patchResponse, patchRequest)
		got := patchResponse.Result().Status
		want := "400 Bad Request"

		if !reflect.DeepEqual(want, got) {
			t.Errorf("\ngot:\n%v\nwant:\n%v\n", got, want)
		}

	})

	t.Run("try to get a bad request in modification again", func(t *testing.T) {

		initDatabase("3")
		defer removeDatabase("3")

		patchRequestBody := strings.NewReader(`{"invalid":"invalid", "invalid":"study", "complete":false}`)
		patchRequest := httptest.NewRequest(http.MethodPatch, "localhost:8080/todolist", patchRequestBody)
		patchResponse := httptest.NewRecorder()

		s.UpdateTodoItem(patchResponse, patchRequest)
		got := patchResponse.Result().Status
		want := "409 Conflict"

		if !reflect.DeepEqual(want, got) {
			t.Errorf("\ngot:\n%v\nwant:\n%v\n", got, want)
		}

	})

}

func TestDeleteTodoItem(t *testing.T) {

	t.Run("try to delete a specific todo item", func(t *testing.T) {

		initDatabase("5")
		defer removeDatabase("5")

		postRequestBody := strings.NewReader(`{"id":1, "todoItem":"study", "complete":false}`)
		postRequest := httptest.NewRequest(http.MethodPost, "localhost:8080/todolist", postRequestBody)
		postResponse := httptest.NewRecorder()

		deleteRequestBody := strings.NewReader(`{"id":1, "todoItem":"study", "complete":false}`)
		deleteRequest := httptest.NewRequest(http.MethodDelete, "localhost:8080/todolist/1", deleteRequestBody)
		deleteResponse := httptest.NewRecorder()

		getRequest := httptest.NewRequest(http.MethodPost, "localhost:8080/todolist", nil)
		getResponse := httptest.NewRecorder()

		want := []Todo{}

		s.AddTodoItem(postResponse, postRequest)
		s.DeleteTodoItem(deleteResponse, deleteRequest)
		s.GetTodoItem(getResponse, getRequest)

		got := []Todo{}
		json.NewDecoder(getResponse.Body).Decode(&got)

		if !reflect.DeepEqual(want, got) {
			t.Errorf("\ngot:\n%v\nwant:\n%v\n", got, want)
		}

	})

	t.Run("try to get a bad request in modification again", func(t *testing.T) {

		initDatabase("3")
		defer removeDatabase("3")

		deleteRequestBody := strings.NewReader(`{"invalid":"invalid", "invalid":"study", "complete":false}`)
		deleteRequest := httptest.NewRequest(http.MethodPatch, "localhost:8080/todolist", deleteRequestBody)
		deleteResponse := httptest.NewRecorder()

		s.DeleteTodoItem(deleteResponse, deleteRequest)
		got := deleteResponse.Result().Status
		want := "400 Bad Request"

		if !reflect.DeepEqual(want, got) {
			t.Errorf("\ngot:\n%v\nwant:\n%v\n", got, want)
		}

	})

}
