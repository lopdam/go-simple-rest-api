package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	// external packages
	"github.com/gorilla/mux"
)

type Task struct {
	Id    int    `json:"id"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

type Tasks []Task

var TasksData = Tasks{
	{
		Id:    1,
		Title: "Example Title",
		Body:  "Example Body",
	},
}

func indexRoute(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Wecome to my GO API!")
}

func createTask(w http.ResponseWriter, r *http.Request) {
	var newTask Task

	reqBody, err := io.ReadAll(r.Body)

	if err != nil {
		fmt.Fprintf(w, "Insert a Valid Task Data")
	}

	json.Unmarshal(reqBody, &newTask)
	newTask.Id = len(TasksData) + 1
	TasksData = append(TasksData, newTask)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTask)

}

func getTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(TasksData)
}

func getOneTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskID, err := strconv.Atoi(vars["id"])
	if err != nil {
		return
	}

	for _, task := range TasksData {
		if task.Id == taskID {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(task)
		}
	}
}

func updateTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskID, err := strconv.Atoi(vars["id"])
	var updatedTask Task

	if err != nil {
		fmt.Fprintf(w, "Invalid ID")
	}

	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Please Enter Valid Data")
	}
	json.Unmarshal(reqBody, &updatedTask)

	for i, t := range TasksData {
		if t.Id == taskID {
			TasksData = append(TasksData[:i], TasksData[i+1:]...)

			updatedTask.Id = t.Id
			TasksData = append(TasksData, updatedTask)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(updatedTask)
			//fmt.Fprintf(w, "The task with ID %v has been updated successfully", taskID)
		}
	}

}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskId, err := strconv.Atoi(vars["id"])

	if err != nil {
		fmt.Fprintf(w, "Invalid User ID")
		return
	}

	for i, t := range TasksData {
		if t.Id == taskId {
			TasksData = append(TasksData[:i], TasksData[i+1:]...)
			fmt.Fprintf(w, "The task with ID %v has been remove successfully", taskId)
		}
	}
}

func main() {
	router := mux.NewRouter().StrictSlash(true)

	// Index Routes
	router.HandleFunc("/", indexRoute)

	// Tasks Routes
	router.HandleFunc("/tasks", createTask).Methods("POST")
	router.HandleFunc("/tasks", getTasks).Methods("GET")
	router.HandleFunc("/tasks/{id}", getOneTask).Methods("GET")
	router.HandleFunc("/tasks/{id}", deleteTask).Methods("DELETE")
	router.HandleFunc("/tasks/{id}", updateTask).Methods("PATCH")

	fmt.Println("Server started on port ", 3000)

	log.Fatal(http.ListenAndServe(":3000", router))
}
