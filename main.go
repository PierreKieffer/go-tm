package main

import (
	"github.com/PierreKieffer/go-tm/pkg/handlers"
	"github.com/PierreKieffer/go-tm/pkg/tm"
	"net/http"
)

func main() {

	taskManager := tm.InitTaskManager()

	http.HandleFunc("/task-manager/trigger", handlers.TaskManagerHandler(taskManager))
	http.HandleFunc("/task-manager/status", handlers.TaskManagerHandler(taskManager))

	fs := http.FileServer(http.Dir("db"))
	corsFS := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		fs.ServeHTTP(w, r)
	})

	http.Handle("/", corsFS)

	http.ListenAndServe(":8080", nil)
}
