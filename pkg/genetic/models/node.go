package models

import "time"

type Node struct {
	ID                string
	AvailableResource Resource
	RemainingResource *Resource
	// On-demand to pay price (￥/h), package in mouthly price is 0
	Price float64
	// Use UTC time to record the node started timestamp
	RunFrom time.Time
	Pods    map[string]Pod

	CpuWeight      float64
	MemoryWeight   float64
	CpuQuotient    float64
	MemoryQuotient float64
}

// Construct a `Node` use all resources, tasks list
func NodeFromPods(id string, all Resource, pods map[string]Pod) Node {
	var remaining *Resource = all.ClonePtr()
	for _, pod := range pods {
		remaining.Sub(pod.RequiredResource)
		(&pod).SetNode(id)
	}
	return Node{
		ID:                id,
		AvailableResource: all,
		RemainingResource: remaining,
		Pods:              pods,
		CpuWeight:         0.,
		MemoryWeight:      0.,
		CpuQuotient:       0.,
		MemoryQuotient:    0.,
	}
}

// Add one Pod to Node
func (n *Node) AddPod(p *Pod) {
	if len(n.Pods) == 0 {
		n.Pods = make(map[string]Pod)
	}
	n.RemainingResource.Sub(p.RequiredResource)
	n.Pods[p.PodID] = *p
	p.SetNode(n.ID)
}
