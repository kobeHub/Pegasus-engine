package models

type Node struct {
	ID                 string
	AvailableResource Resource
	RemainingResource *Resource
	Pods              map[string]Pod

	CpuWeight      float64
	MemoryWeight   float64
	CpuQuotient    float64
	MemoryQuotient float64
}

// Construct a `Node` use all resources, tasks list
func NewNode(id string, all Resource, pods map[string]Pod) Node {
	var remaining *Resource = all.ClonePtr()
	for _, pod := range pods {
		remaining.Sub(pod.RequiredResource)
		pod.SetNode(id)
	}
	return Node{
		ID:                 id,
		AvailableResource: all,
		RemainingResource: remaining,
		Pods:              pods,
		CpuWeight:          0.,
		MemoryWeight:       0.,
		CpuQuotient:        0.,
		MemoryQuotient:     0.,
	}
}
