package main

import (
	"net/http"

	"github.com/Boalq/worker-colony/internal/handler"
)

func main() {
	tsk := handler.NewTaskHandler()

	// Set the amount of worker (bees) for concurrent work!
	for i := 1; i <= 3; i++ {
        go tsk.StartWorker(i)
    }

	mux := http.NewServeMux()

	//Handing a task to the http server
	mux.HandleFunc("POST /task", tsk.Tasks)

	// Checking on the state of the Task
	mux.HandleFunc("GET /task/{id}", tsk.GetTask)

	http.ListenAndServe("localhost:8080", mux)
}
