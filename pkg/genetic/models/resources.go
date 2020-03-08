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

func (info *Resource) Less(ano Resource) bool {
	return info.cpuCores <= ano.cpuCores && info.memory <= ano.memory
}

func (info *Resource) CpuRest(avail Resource) float64 {
	return info.cpuCores / avail.cpuCores
}

func (info *Resource) MemRest(avail Resource) float64 {
	return info.memory / avail.memory
}

func (info *Resource) NotAvail() bool {
	return info.cpuCores < 0 || info.memory < 0
}

func (info *Resource) NoMem() bool {
	return info.memory < 0
}

func (info *Resource) NoCpu() bool {
	return info.cpuCores < 0
}

func (info *Resource) Mem() float64 {
	return info.memory
}

func (info *Resource) Cpu() float64 {
	return info.cpuCores
}
