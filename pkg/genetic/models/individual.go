/*
`models` package defaines the nessary struct to implements
genetic algorithm. including:
   + Individual
   + Task
   + Node
   + Resources
*/
package models

import "math"

// The item of genetic algorithm, one individual is
// is one entry during the evolution
type Individual struct {
	ID string
	// The taskId as key and nodeId as value, one `Assignment`
	// is a schedule scheme
	Assignment map[string]string
	// The pods layout before rescheduling
	OrignalAssignment    map[string]string
	AllNodes             map[string]Node
	AllPods             map[string]Pod
	NumOfUnassignedPods int
	NumOfUnassignedNodes int

	// Objective valus, all the objective values to be optimized
	ObjectiveValues           []float64
	TranslatedObjectiveValues []float64
	NormalizedObjectiveValues []float64
	PerpendicularDistance     float64
	ReferencePoint            ReferencePoint

	// NSGA-II
	IndividualsDominatedByThis   []*Individual
	NumOfIndividualsDominateThis int
	Rank                         int
	CrowdingDistance             float64

	// NSGA-III
	ConstrainedViolationValue float64
	IsFeasible                bool
}

// The initial state of current cluster before rescheduling
// init the individual
//
// every `Node` contains tasks list and resources usage state
// `num_tasks` is the number of current pods
func (info *Individual) init(nodes []Node, pending_pods []Pod, num_run_pods int) {
	info.OrignalAssignment = make(map[string]string, num_run_pods)
	info.AllPods = make(map[string]Pod, num_run_pods+len(pending_pods))
	info.AllNodes = make(map[string]Node, len(nodes))
	info.NumOfUnassignedPods = len(pending_pods)

	for _, pt := range pending_pods {
		info.AllPods[pt.PodID] = pt
	}
	for _, node := range nodes {
		info.AllNodes[node.ID] = node
		if len(node.Pods) == 0 {
			info.NumOfUnassignedNodes += 1
			continue
		}
		for key, pod := range node.Pods {
			info.AllPods[key] = pod
			info.OrignalAssignment[pod.PodID] = node.ID
		}
	}
	info.computeObjectiveValues()
}

func (info *Individual) computeObjectiveValues() {
	info.ObjectiveValues = make([]float64, 3)
	info.ObjectiveValues = append(info.ObjectiveValues, info.spreadObjective())
	info.ObjectiveValues = append(info.ObjectiveValues, info.uniquenessObjective())
	info.ObjectiveValues = append(info.ObjectiveValues, info.resourceUtilization())
}

// Compute the pods distribution value, the smaller has
// higher distrubutivity
func (info *Individual) spreadObjective() float64 {
	spreadValue := 0

	for _, node := range info.AllNodes {
		nodeSpread := 0
		for i := 0; i < len(node.Pods); i++ {
			nodeSpread += i + 1
		}
		spreadValue += nodeSpread
	}

	return float64(spreadValue)
}

// Conpute the Pod Uniqueness of every node
func (info *Individual) uniquenessObjective() float64 {
	value := 0

	for _, node := range info.AllNodes {
		type_cnt := make(map[string]int)
		for _, pod := range node.Pods {
			type_cnt[pod.PodType] += 1
		}
		nodeUnique := 0
		for _, cnt := range type_cnt {
			for i := 0; i < cnt; i++ {
				nodeUnique += i + 1
			}
		}
		value += nodeUnique
	}
	return float64(value)
}

func (info *Individual) resourceUtilization() float64 {
	value := 0.
	for _, node := range info.AllNodes {
		value += math.Abs(
			node.RemainingResource.CpuRest(node.AvailableResource) -
				node.RemainingResource.MemRest(node.AvailableResource))
	}
	return value
}
