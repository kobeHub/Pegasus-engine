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
	PodID            string   `json:"podId,required"`
	RequiredResource Resource `json:"resource,required"`
	NodeID           string   `json:"nodeId,omitempty"`
	Status           PodPhase `json:"status"`
	PodType          string   `json:"podType,omitempty"`
}

func NewPod(id string, cpu, mem float64, nid string, status PodPhase) Pod {
	return Pod{
		PodID:            id,
		RequiredResource: NewResource(cpu, mem),
		NodeID:           nid,
		Status:           status,
		PodType:          "",
	}
}

func (info *Pod) SetNode(id string) {
	info.NodeID = id
}
