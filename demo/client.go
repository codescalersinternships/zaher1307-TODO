package main

import (
	"bufio"
	"bytes"
	"fmt"
	"net/http"
)

const (
    Req = `{"id": %d, "todoItem": "task ID: %d", "complete": %s}`
    ServerAddress = "http://localhost:8080%s"

    MethodGet     = "GET"
    MethodHead    = "HEAD"
    MethodPost    = "POST"
    MethodPut     = "PUT"
    MethodPatch   = "PATCH"
    MethodDelete  = "DELETE"
    MethodConnect = "CONNECT"
    MethodOptions = "OPTIONS"
    MethodTrace   = "TRACE"
)

func main() {
    client := &http.Client{}

    // Adding new task
    fmt.Println("*********Add new tasks*********")
    for i := 1; i <= 10; i++ {
        body := fmt.Sprintf(Req, i, i, "false")
        fmt.Println(body)
        req := buildRequest(MethodPost, "/todolist", []byte(body))
        res, _ := client.Do(req)
        printResponse(res)
    }

    // Get all tasks
    fmt.Println("*********Get all tasks*********")
    req := buildRequest(MethodGet, "/todolist", []byte(""))
    res, _ := client.Do(req)
    printResponse(res)

    // Get task by ID
    fmt.Println("*********Get task by ID = 1*********")
    req = buildRequest(MethodGet, "/todolist/1", []byte(""))
    res, _ = client.Do(req)
    printResponse(res)

    // Update task by ID
    fmt.Println("*********Update task with ID = 1*********")
    body := fmt.Sprintf(Req, 1, 1, "true")
    req = buildRequest(MethodPatch, "/todolist", []byte(body))
    res, _ = client.Do(req)
    printResponse(res)

    // Delete task by ID
    fmt.Println("*********Delete task by ID = 1*********")
    req = buildRequest(MethodDelete, "/todolist/1", []byte(""))
    res, _ = client.Do(req)
    printResponse(res)

    fmt.Println("*********Delete task by ID = 2*********")
    req = buildRequest(MethodDelete, "/todolist/2", []byte(""))
    res, _ = client.Do(req)
    printResponse(res)

    fmt.Println("*********Delete task by ID = 3*********")
    req = buildRequest(MethodDelete, "/todolist/3", []byte(""))
    res, _ = client.Do(req)
    printResponse(res)

    // Get all tasks
    fmt.Println("*********Get all tasks*********")
    req = buildRequest(MethodGet, "/todolist", []byte(""))
    res, _ = client.Do(req)
    printResponse(res)
}

func buildRequest(method, endpoint string, body []byte) *http.Request {
    serverAddress := fmt.Sprintf(ServerAddress, endpoint)
    req, _ := http.NewRequest(method, serverAddress, bytes.NewBuffer(body))
    req.Header.Set("Content-Type", "application/json")

    return req
}
func printResponse(res *http.Response) {
    defer res.Body.Close()
    fmt.Println("Response status:", res.Status)
    scanner := bufio.NewScanner(res.Body)
    for scanner.Scan() {
        fmt.Println(scanner.Text())
    }
    fmt.Println()
}
