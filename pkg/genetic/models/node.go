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
