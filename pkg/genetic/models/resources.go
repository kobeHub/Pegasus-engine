package models

type Resource struct {
	CpuCores float64 `json:"cpu"`
	Memory   float64 `json:"memory"`
}

func NewResource(cpuCores, memory float64) Resource {
	return Resource{
		CpuCores: cpuCores,
		Memory:   memory,
	}
}

func (info *Resource) Sub(ano Resource) {
	info.CpuCores -= ano.CpuCores
	info.Memory -= ano.Memory
}

func (info *Resource) ClonePtr() *Resource {
	return &Resource{CpuCores: info.CpuCores, Memory: info.Memory}
}

func (info *Resource) Less(ano Resource) bool {
	return (info.CpuCores <= ano.CpuCores) && (info.Memory <= ano.Memory)
}

func (info *Resource) CpuRest(avail Resource) float64 {
	return info.CpuCores / avail.CpuCores
}

func (info *Resource) MemRest(avail Resource) float64 {
	return info.Memory / avail.Memory
}

func (info *Resource) NotAvail() bool {
	return info.CpuCores < 0 || info.Memory < 0
}

func (info *Resource) NoMem() bool {
	return info.Memory < 0
}

func (info *Resource) NoCpu() bool {
	return info.CpuCores < 0
}

func (info *Resource) Mem() float64 {
	return info.Memory
}

func (info *Resource) Cpu() float64 {
	return info.CpuCores
}
