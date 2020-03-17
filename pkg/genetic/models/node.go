package models

import "time"

type Node struct {
	ID                string    `json:"id,required"`
	Name              string    `json:"name,required"`
	AvailableResource Resource  `json:"availableResource,required"`
	RemainingResource *Resource `json:"remainingResouece,omitempty"`
	// On-demand to pay price (￥/h), package in mouthly price is 0
	Price float64 `json:"price,omitempty"`
	// Use UTC time to record the node started timestamp
	RunFrom time.Time      `json:"runFrom,omitempty"`
	Pods    map[string]Pod `json:"pods,omitempty"`

	CpuWeight      float64 `json:"-"`
	MemoryWeight   float64 `json:"-"`
	CpuQuotient    float64 `json:"-"`
	MemoryQuotient float64 `json:"-"`
}

// Construct a `Node` use all resources, tasks list
func NewConsistNode(id, name string, all Resource, runFrom time.Time) Node {
	var remaining *Resource = all.ClonePtr()
	return Node{
		ID:                id,
		Name:              name,
		AvailableResource: all,
		RemainingResource: remaining,
		Price:             0.,
		RunFrom:           runFrom,
		CpuWeight:         0.,
		MemoryWeight:      0.,
		CpuQuotient:       0.,
		MemoryQuotient:    0.,
	}
}

func NewDemandNode(id, name string, all Resource, price float64, run_from time.Time) Node {
	var remaining *Resource = all.ClonePtr()
	return Node{
		ID:                id,
		Name:              name,
		AvailableResource: all,
		RemainingResource: remaining,
		Price:             price,
		RunFrom:           run_from,
		CpuWeight:         0.,
		MemoryWeight:      0.,
		CpuQuotient:       0.,
		MemoryQuotient:    0.,
	}
}

// Add one Pod to Node
func (n *Node) AddPod(p *Pod) {
	if n.Pods == nil {
		n.Pods = make(map[string]Pod)
	}
	n.RemainingResource.Sub(p.RequiredResource)
	p.SetNode(n.ID)
	n.Pods[p.PodID] = *p
}
