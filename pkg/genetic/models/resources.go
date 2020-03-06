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

func (info *Resource) Sub(ano Resource) {
	info.cpuCores -= ano.cpuCores
	info.memory -= ano.memory
}

func (info *Resource) ClonePtr() *Resource {
	return &Resource{info.cpuCores, info.memory}
}

func (info *Resource) CpuRest(avail Resource) float64 {
	return info.cpuCores / avail.cpuCores
}

func (info *Resource) MemRest(avail Resource) float64 {
	return info.memory / avail.memory
}
