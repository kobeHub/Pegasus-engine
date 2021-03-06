/*
`models` package defaines the nessary struct to implements
genetic algorithm. including:
   + Individual
   + Task
   + Node
   + Resources
*/
package models

import (
	"math"
	_ "time"

	log "github.com/sirupsen/logrus"
)

// The item of genetic algorithm, one individual is
// is one entry during the evolution
type Individual struct {
	ID string
	// The taskId as key and nodeId as value, one `Assignment`
	// is a schedule scheme
	Assignment map[string]string
	// The pods layout before rescheduling
	OriginalAssignment   map[string]string
	AllNodes             map[string]Node
	AllPods              map[string]Pod
	NumOfUnassignedPods  int
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
	Rank                         int // smalller superior
	CrowdingDistance             float64

	// NSGA-III
	ConstraintedViolationValue float64
	IsFeasible                 bool
}

// The `Population` of one generation, includes of individuals
type Population []*Individual

// Ordered individuals
type Front []*Individual
type Fronts []*Front

// Initialize Individual according to `assign`, and all nodes
// and Pod resource is known, set `Node` and `Pod` field
//
// Before `Init` func, the Individual's Assignment(and original) need to be
// set properly
//
// Node in `nodes` known: ID, AvailableResource, RemainingResource
// Pod in `pods` known: PodID, Status, RequiredResource
func (info *Individual) Init(nodes []Node, pods []Pod) {
	info.AllPods = make(map[string]Pod, len(pods))
	info.AllNodes = make(map[string]Node, len(nodes))

	for _, node := range nodes {
		// reset node's pods list
		for pod := range node.Pods {
			delete(node.Pods, pod)
		}
		node.RemainingResource = node.AvailableResource.ClonePtr()
		info.AllNodes[node.ID] = node
	}
	for _, pod := range pods {
		info.AllPods[pod.PodID] = pod
	}

	info.NumOfUnassignedPods = 0
	info.NumOfUnassignedNodes = 0

	for pid, nid := range info.Assignment {
		if len(nid) == 0 {
			info.NumOfUnassignedPods += 1
			continue
		}

		pod := info.AllPods[pid]
		node := info.AllNodes[nid]
		log.WithFields(log.Fields{
			"node-remain": node.RemainingResource,
			"pod-require": pod.RequiredResource,
		}).Debug("Add pod before")

		node.AddPod(&pod)

		log.WithFields(log.Fields{
			"node-remain": node.RemainingResource,
			"pod-require": pod.RequiredResource,
		}).Debug("Add pod after")

		info.AllPods[pid] = pod
		info.AllNodes[nid] = node
	}

	for _, node := range info.AllNodes {
		if len(node.Pods) == 0 {
			info.NumOfUnassignedNodes += 1
		}
	}

	info.ComputeValues()
}

func (info *Individual) ComputeValues() {
	info.ObjectiveValues = make([]float64, 0, 4)
	info.TranslatedObjectiveValues = make([]float64, 4)
	info.NormalizedObjectiveValues = make([]float64, 4)
	info.ObjectiveValues = append(info.ObjectiveValues, info.ondemandPrice())
	info.ObjectiveValues = append(info.ObjectiveValues, info.spreadObjective())
	info.ObjectiveValues = append(info.ObjectiveValues, info.uniquenessObjective())
	info.ObjectiveValues = append(info.ObjectiveValues, info.resourceUtilization())
	info.ComputeIsFeasible()
	info.ComputeViolationValue()
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

func (info *Individual) ondemandPrice() float64 {
	price := 0.
	for _, node := range info.AllNodes {
		if node.IsDemand && len(node.Pods) != 0 {
			price +=  24 * node.Price
		}
	}
	return price
}

// Raw check this `Individual` is superior than ano or not
func (info *Individual) dominates(ano Individual) bool {
	// no worse than
	for i := 0; i < len(info.ObjectiveValues); i++ {
		if info.ObjectiveValues[i] > ano.ObjectiveValues[i] {
			return false
		}
	}

	// One value is better than another's
	for i := 0; i < len(info.ObjectiveValues); i++ {
		if info.ObjectiveValues[i] < ano.ObjectiveValues[i] {
			return true
		}
	}

	return false
}

// crowdwd compate operator
func (info *Individual) CrowdedCompareLess(ano Individual) bool {
	return info.Rank < ano.Rank ||
		((info.Rank == ano.Rank) && (info.CrowdingDistance > ano.CrowdingDistance))
}

// ********************* Feasible and constrained ***************************

func (info *Individual) ComputeIsFeasible() {
	for _, node := range info.AllNodes {
		if node.RemainingResource.NotAvail() {
			info.IsFeasible = false
			return
		}
	}

	info.IsFeasible = true
}

func (info *Individual) ComputeViolationValue() {
	value := 0.

	for _, node := range info.AllNodes {
		if node.RemainingResource.NoMem() {
			value += math.Abs(node.RemainingResource.Mem())
		}

		if node.RemainingResource.NoCpu() {
			value += math.Abs(node.RemainingResource.Cpu())
		}
	}
	info.ConstraintedViolationValue = value
}

func (info *Individual) ConstraintDominate(ano Individual) bool {
	if (info.IsFeasible && !ano.IsFeasible) ||
		(!info.IsFeasible && !ano.IsFeasible && info.ConstraintedViolationValue < ano.ConstraintedViolationValue) || (info.IsFeasible && ano.IsFeasible && info.dominates(ano)) {
		return true
	} else {
		return false
	}
}

/// Population best price Individual, return false means all the individual is
/// not feasible
func (popu Population) BestPriceIndividual() (Individual, float64, bool) {
	bestPrice := math.MaxFloat64
	var bestId int
	for id, info := range popu {
		if info.IsFeasible && info.ObjectiveValues[0] < bestPrice {
			bestPrice = info.ObjectiveValues[0]
			bestId = id
		}
	}
	return *popu[bestId], bestPrice, bestPrice != math.MaxFloat64
}
