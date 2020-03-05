package models

type Task struct {
	TaskID           string
	RequiredResource Resource
	NodeID           string
	TaskType         string
}

func NewTask(id string, cpu, mem float64) Task {
	return Task{
		TaskID:           id,
		RequiredResource: NewResource(cpu, mem),
		NodeID:           "",
		TaskType:         "",
	}
}
