package models

type Resource struct {
	cpuCores float64
	memory   float64
}

func NewResource(cpuCores, memory float64) Resource {
	return Resource{
		cpuCores,
		memory,
	}
}

func (info *Resource) Sub(ano *Resource) {
	info.cpuCores -= ano.cpuCores
	info.memory -= ano.memory
}
