package main

import (
	"net/http"

	"github.com/Boalq/worker-colony/internal/handler"
)

func main() {
	tsk := handler.NewTaskHandler()
	mux := http.NewServeMux()
	mux.HandleFunc("/task", tsk.Tasks)

	http.ListenAndServe("localhost:8080", mux)
}
