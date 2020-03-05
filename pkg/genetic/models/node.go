package models

type Node struct {
	ID                 string
	AvaliableResources Resource
	RemainingResources *Resource
	Tasks              map[string]Task

	CpuWeight      float64
	MemoryWeight   float64
	CpuQuotient    float64
	MemoryQuotient float64
}

// Construct a `Node` use all resources, tasks list
func NewNode(id string, all Resource, tasks map[string]Task) Node {
	var remaining *Resource = all.ClonePtr()
	for _, task := range tasks {
		remaining.Sub(task.RequiredResource)
	}
	return Node{
		ID:                 id,
		AvaliableResources: all,
		RemainingResources: remaining,
		Tasks:              tasks,
		CpuWeight:          0.,
		MemoryWeight:       0.,
		CpuQuotient:        0.,
		MemoryQuotient:     0.,
	}
}
