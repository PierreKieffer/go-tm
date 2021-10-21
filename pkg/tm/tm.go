package tm

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"sync"
	"time"
)

type TaskManager struct {
	mu       sync.Mutex
	Tasks    map[string]Task `json:"tasks"`
	Database string          `json:"database"`
}

type Task struct {
	TaskID   string                 `json:"taskID"`
	TaskType string                 `json:"taskType"`
	Status   string                 `json:"status"`
	Meta     map[string]interface{} `json:"meta"`
}

func InitTaskManager() *TaskManager {
	/*
		Init task manager
	*/

	fmt.Println("Init task manager ...")

	var taskManager = &TaskManager{}
	taskManager.Database = "task-manager-db.json"

	_, err := os.Stat(taskManager.Database)
	if err != nil {
		taskManager.Tasks = make(map[string]Task)
		fmt.Println("task manager is ready")
		return taskManager
	}
	database, _ := os.Open(taskManager.Database)
	defer database.Close()

	buffer, _ := ioutil.ReadAll(database)
	var tasks map[string]Task
	json.Unmarshal(buffer, &tasks)

	taskManager.Tasks = tasks

	fmt.Println("task manager is ready")
	return taskManager
}

func (taskManager *TaskManager) CreateTask(taskType string) Task {
	/*
		Create a new task in task manager
	*/

	taskManager.mu.Lock()
	defer taskManager.mu.Unlock()

	// Create task
	taskID := strconv.FormatInt(time.Now().UnixNano(), 10)

	task := Task{
		TaskID:   taskID,
		TaskType: taskType,
		Status:   "PENDING",
	}

	taskManager.Tasks[task.TaskID] = task

	taskJson, _ := json.Marshal(task)

	// Save task metadata
	taskManager.SaveTaskManager()

	fmt.Printf("Task created : %v\n", string(taskJson))

	return task
}

func (taskManager *TaskManager) SaveTaskManager() {
	/*
		Save state of task manager
	*/

	buffer, _ := json.Marshal(taskManager.Tasks)
	ioutil.WriteFile(taskManager.Database, buffer, 0755)
}

func (taskManager *TaskManager) DeleteTask(taskID string) {
	/*
		Delete task if status if done
	*/

	_, ok := taskManager.Tasks[taskID]
	if ok {
		if taskManager.Tasks[taskID].Status == "DONE" {
			delete(taskManager.Tasks, taskID)
		}
	}

	taskManager.SaveTaskManager()
}

func (taskManager *TaskManager) Task(taskID string) Task {
	/*
		Return task
	*/

	taskManager.mu.Lock()
	defer taskManager.mu.Unlock()

	task := taskManager.Tasks[taskID]
	return task
}

func (taskManager *TaskManager) UpdateTaskStatus(taskID string, status string) {
	/*
		Update task status
	*/

	taskManager.mu.Lock()
	defer taskManager.mu.Unlock()

	task, ok := taskManager.Tasks[taskID]
	if ok {
		task.Status = status
		taskManager.Tasks[taskID] = task

		taskManager.SaveTaskManager()
	}
}

func (taskManager *TaskManager) UpdateTaskMeta(taskID string, key string, value interface{}) {
	/*
	   Update task metadata
	*/

	taskManager.mu.Lock()
	defer taskManager.mu.Unlock()

	task, ok := taskManager.Tasks[taskID]
	if ok {
		task.Meta[key] = value
		taskManager.Tasks[taskID] = task

		taskManager.SaveTaskManager()
	}
}
