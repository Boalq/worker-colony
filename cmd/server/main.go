package main

import (
	"net/http"

	"github.com/Boalq/worker-colony/internal/handler"
)

func main() {
	tsk := handler.NewTaskHandler()

	for i := 1; i <= 3; i++ {
        go tsk.StartWorker(i)
    }

	mux := http.NewServeMux()
	mux.HandleFunc("POST /task", tsk.Tasks)
	mux.HandleFunc("GET /task/{id}", tsk.GetTask)

	http.ListenAndServe("localhost:8080", mux)
}
