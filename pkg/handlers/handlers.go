package handlers

import (
	"encoding/json"
	"github.com/PierreKieffer/go-tm/pkg/tm"
	"log"
	"net/http"
)

func TaskManagerHandler(taskManager *tm.TaskManager) http.HandlerFunc {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		switch request.Method {
		case http.MethodPost:

			switch request.URL.Path {
			case "/task-manager/trigger":
				payloadBuffer := json.NewDecoder(request.Body)
				var payload map[string]interface{}
				payloadBuffer.Decode(&payload)

				jsonPayload, _ := json.Marshal(payload)

				log.Printf("[POST] /task-manager/trigger : %v", string(jsonPayload))

				// Create task
				task := taskManager.CreateTask(payload["taskType"].(string))

				// Launch task
				taskManager.UpdateTaskStatus(task.TaskID, "RUNNING")

				// Response
				response, _ := json.Marshal(task)
				writer.Header().Add("Content-type", "application/json")
				writer.Write(response)
			}

		case http.MethodGet:
			switch request.URL.Path {
			case "/task-manager/status":
				taskID := request.URL.Query()["taskID"][0]
				log.Printf("[GET] /task-manager/status : %v", taskID)
				task := taskManager.Task(taskID)

				response, _ := json.Marshal(task)
				writer.Header().Add("Content-type", "application/json")
				writer.Write(response)
			}
		}
	})
}
