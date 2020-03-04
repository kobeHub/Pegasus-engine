package models

type Task struct {
	TaskID            string
	RequiredResources Resource
	NodeID            string
	TaskType          string
}
