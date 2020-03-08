package models

// The phase of kubernetes Pod
type PodPhase string

const (
	Pending   PodPhase = "Pending"
	Runing    PodPhase = "Runing"
	Succeeded PodPhase = "Succeeded"
	Failed    PodPhase = "Failed"
	Unknown   PodPhase = "Unknown"
)

type Pod struct {
	PodID            string
	RequiredResource Resource
	NodeID           string
	Status           PodPhase
	PodType          string
}

func NewPod(id string, cpu, mem float64) Pod {
	return Pod{
		PodID:            id,
		RequiredResource: NewResource(cpu, mem),
		NodeID:           "",
		Status:           Unknown,
		PodType:          "",
	}
}

func (info *Pod) SetNode(id string) {
	info.NodeID = id
}
