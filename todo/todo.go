package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

func checkError(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
}

// get all the tasks from the json file and parse into list of tasks
func getTasks() (tasks []Task) {
	taskBytes, err := os.ReadFile("./testdata/tasks.json")
	checkError(err)

	err = json.Unmarshal(taskBytes, &tasks)
	checkError(err)

	return tasks
}

// save books to books.json file
func saveTasks(tasks []Task) error {

	// convert into bytes for writing to the JSON file
	taskBytes, err := json.Marshal(tasks)
	checkError(err)

	err = os.WriteFile("./testdata/tasks.json", taskBytes, 0644)

	return err
}

// get tasks logic
func handleGetTasks(getCmd *flag.FlagSet, all *bool, id *string) {

	getCmd.Parse(os.Args[2:])

	// check for all or by id
	if !*all && *id == "" {
		fmt.Println("subcommand --all or --id needed")
		getCmd.PrintDefaults()
		os.Exit(1)
	}

	// if for all return all books
	if *all {
		tasks := getTasks()
		fmt.Printf("Id \t Title \t isCompleted \t createdAt \n")

		for _, task := range tasks {
			fmt.Printf("%v \t %v \t %v \t %v \n", task.Id, task.Title, task.IsCompleted, task.CreatedAt)
		}
	}

	// if by id, return only that task
	// throw error if task not found
	if *id != "" {
		tasks := getTasks()
		fmt.Printf("Id \t Title \t isCompleted \t createdAt \n")
		// check if book exists
		var foundTask bool
		for _, task := range tasks {
			if *id == task.Id {
				foundTask = true
				fmt.Printf("%v \t %v \t %v \t %v \n", task.Id, task.Title, task.IsCompleted, task.CreatedAt)
			}
		}
		// if no task found with this id, throw an error
		if !foundTask {
			fmt.Println("Task not found")
			os.Exit(1)
		}
	}
}

// logic for add or update a task
func handleAddTask(addCmd *flag.FlagSet, id, title *string, isCompleted *bool, createdAt *string, addNewTask bool) {
	addCmd.Parse(os.Args[2:])

	if *id == "" || *title == "" || *createdAt == "" {
		fmt.Println("Please provide task id, title, createdAt")
		addCmd.PrintDefaults()
		os.Exit(1)
	}

	tasks := getTasks()
	var newTask Task
	// check if a task exists
	var foundTask bool

	// checking whether to add or to update
	if addNewTask {
		newTask = Task{*id, *title, *isCompleted, *createdAt}
		tasks = append(tasks, newTask)
	} else {
		for i, task := range tasks {
			if task.Id == *id {
				// replace old values with new ones
				tasks[i] = Task{*id, *title, *isCompleted, *createdAt}
				foundTask = true
			}
		}

		// if no task found with id, throw an error
		if !foundTask {
			fmt.Println("Task not found")
			os.Exit(1)
		}
	}

	err := saveTasks(tasks)
	checkError(err)

	fmt.Println("Task added successfully")
}

func handleDeleteTask(deleteCmd *flag.FlagSet, id *string) {
	deleteCmd.Parse(os.Args[2:])

	if *id == "" {
		fmt.Println("Please provide book --id")
		deleteCmd.PrintDefaults()
		os.Exit(1)
	}

	tasks := getTasks()
	var foundTask bool

	for i, task := range tasks {
		if task.Id == *id {
			//delete from list
			tasks = append(tasks[:i], tasks[i+1:]...)
			foundTask = true
		}
	}

	// if no task found with id, throw an error
	if !foundTask {
		fmt.Println("Task not found")
		os.Exit(1)
	}

	err := saveTasks(tasks)
	checkError(err)

	fmt.Println("Task deleted successfully")
}
