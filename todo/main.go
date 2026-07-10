package main

import (
	"flag"
	"fmt"
	"os"
)

// struct based on tasks.json file
type Task struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	IsCompleted bool   `json:"is_completed"`
	CreatedAt   string `json:"created_at"`
}

func main() {
	/*
		get tasks --all or --id
		./todo get --all
		./todo get --id=4
	*/
	getCmd := flag.NewFlagSet("get", flag.ExitOnError)
	getAll := getCmd.Bool("all", false, "List all the tasks")
	getId := getCmd.String("id", "", "Get task by id")
	/*
		add a task with id, title, isCompleted, and createdAt
		./todo add --id=6 --title=test-task --is_completed=false --created_at=7/9/20268:34PM
	*/
	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	addId := addCmd.String("id", "", "Task id")
	addTitle := addCmd.String("title", "", "Task title")
	addIsCompleted := addCmd.Bool("is_completed", false, "Whether task is completed")
	addCreatedAt := addCmd.String("created_at", "", "Date and time we set the task for ourselves")
	/*
		update a task with id, title, is_completed, and created_at
		./todo update --id=6 --title=test-task-1 --is_completed=true --created_at=7/9/20268:45PM
	*/
	updateCmd := flag.NewFlagSet("update", flag.ExitOnError)
	updateId := updateCmd.String("id", "", "Task id")
	updateTitle := updateCmd.String("title", "", "Task title")
	updateIsCompleted := updateCmd.Bool("is_completed", false, "Whether the task is completed")
	updateCreatedAt := updateCmd.String("created_at", "", "Date and time we set the task for ourselves")
	/*
		delete a book by --id
		./todo delete --id=6
	*/
	deleteCmd := flag.NewFlagSet("delete", flag.ExitOnError)
	deleteId := deleteCmd.String("id", "", "Delete task by id")

	// validate the command line arguments
	if len(os.Args) < 2 {
		fmt.Println("Expected get, add, update, or delete subcommand")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "get":
		handleGetTasks(getCmd, getAll, getId)
	case "add":
		handleAddTask(addCmd, addId, addTitle, addIsCompleted, addCreatedAt, true)
	case "update":
		handleAddTask(updateCmd, updateId, updateTitle, updateIsCompleted, updateCreatedAt, false)
	case "delete":
		handleDeleteTask(deleteCmd, deleteId)
	default:
		fmt.Println("Please provide get, add, update, or delete subcommand")
		os.Exit(1)
	}
}
