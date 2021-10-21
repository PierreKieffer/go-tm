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
	http.ListenAndServe(":8080", nil)
}
