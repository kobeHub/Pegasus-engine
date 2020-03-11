package genetic

import (
	_ "fmt"
	"github.com/rs/xid"
	_ "math"
	"math/rand"
	"time"

	"github.com/kobeHub/Pegasus-engine/pkg/genetic/models"
)

type Genetic struct {
	AllNodes           []models.Node
	AllPods            []models.Pod
	OriginalAssignment map[string]string
	Size               int
	GenerationNum      int
}

// Generate random feasible assign according to current
// Pods and Nodes
func (g Genetic) GenerateRandomFeasibleIndividual() *models.Individual {
	nodes := make([]models.Node, len(g.AllNodes))

	// Reset Nodes remaining resources
	for i, node := range g.AllNodes {
		nodes[i] = models.Node{
			RemainingResource: node.AvailableResource.ClonePtr(),
		}
	}

	shuffleNodes(nodes)

	uid := xid.New()
	assign := make(map[string]string, len(g.AllPods))
	for _, pod := range g.AllPods {
		for _, node := range nodes {
			if pod.RequiredResource.Less(*node.RemainingResource) {
				assign[pod.PodID] = node.ID
				node.RemainingResource.Sub(pod.RequiredResource)
			}
			break
		}
	}

	info := models.Individual{
		ID:                 uid.String(),
		Assignment:         assign,
		OriginalAssignment: g.OriginalAssignment,
	}

	info.Init(g.AllNodes, g.AllPods)
	return &info
}

func (g Genetic) GenerateRandomFeasiblePopulation() models.Population {
	popu := make([]*models.Individual, g.Size)
	for i := 0; i < g.Size; i++ {
		popu[i] = g.GenerateRandomFeasibleIndividual()
	}
	return popu
}

// Use cluster current state initialize population
func (g Genetic) GenerateInitPopulation() models.Population {
	popu := make([]*models.Individual, g.Size)
	return popu
}

// Combine two Populations into one
func (g Genetic) combinePopulation(first, second models.Population) models.Population {
	result := make([]*models.Individual, g.Size*2)
	for i := 0; i < g.Size; i++ {
		result[i] = first[i]
		result[i+g.Size] = second[i]
	}
	return result
}

// Select superior one from random two
func (g Genetic) binarySelect(popu models.Population) models.Individual {
	first := popu[rand.Intn(g.Size)]
	second := popu[rand.Intn(g.Size)]

	if first.CrowdedCompareLess(*second) {
		return *first
	} else {
		return *second
	}
}

// Constrainted nsga iii
func (g Genetic) constraintedBinarySelect(popu models.Population) models.Individual {
	first := popu[rand.Intn(g.Size/2)]
	second := popu[rand.Intn(g.Size/2+rand.Intn(g.Size/2))]

	if first.IsFeasible && !second.IsFeasible {
		return *first
	} else if !first.IsFeasible && second.IsFeasible {
		return *second
	} else if !first.IsFeasible && !second.IsFeasible {
		if first.ConstraintedViolationValue > second.ConstraintedViolationValue {
			return *second
		} else if first.ConstraintedViolationValue < second.ConstraintedViolationValue {
			return *first
		} else {
			if rand.Float64() > 0.5 {
				return *first
			} else {
				return *second
			}
		}
	} else {
		if rand.Float64() > 0.5 {
			return *first
		} else {
			return *second
		}
	}
}

// Make new population from parent
func (g Genetic) makeNewPopulation(parent models.Population) models.Population {
	new := make([]*models.Individual, g.Size)
	var (
		first  models.Individual
		second models.Individual
	)
	for i := 0; i < g.Size; i++ {
		first = g.constraintedBinarySelect(parent)
		second = g.constraintedBinarySelect(parent)
		newIndi := g.reproduce(first, second)

		if rand.Float64() > 0.5 {
			g.mutate(&newIndi)
		}
		new[i] = &newIndi
	}
	return new
}

// Run genetic algirithm NSGA3 implements
func (g Genetic) RunGeneticNSGA3(num_segaments int) models.Population {
	nsga3 := NSGA3{
		Ops: NSGA2{},
	}
	parent := g.GenerateInitPopulation()

	for t := 0; t < g.GenerationNum; t++ {
		rps := nsga3.GetReferencePoints(len(parent[0].ObjectiveValues), num_segaments)
		nextPopu := nsga3.GenerateNextPopulation(t, g, parent, rps)
		parent = nextPopu
		for _, indi := range parent {
			indi.ReferencePoint = models.ReferencePoint{}
			indi.PerpendicularDistance = 0
			indi.IndividualsDominatedByThis = []*models.Individual{}
			indi.NumOfIndividualsDominateThis = 0
			indi.Rank = 0
		}
	}
	return parent
}

//******** Individual operations to generate new *********

// Select random point to crossover
func (g Genetic) reproduce(first, second models.Individual) models.Individual {
	num_pods := len(first.AllPods)
	newAssign := make(map[string]string, num_pods)
	randomCut := rand.Intn(num_pods)

	// Before cutpoint from a, after from b
	for i := 0; i < num_pods; i++ {
		pod := g.AllPods[i]
		if i < randomCut {
			newAssign[pod.PodID] = first.Assignment[pod.PodID]
		} else {
			newAssign[pod.PodID] = second.Assignment[pod.PodID]
		}
	}

	uid := xid.New()
	info := models.Individual{
		ID:                 uid.String(),
		Assignment:         newAssign,
		OriginalAssignment: g.OriginalAssignment,
	}
	info.Init(g.AllNodes, g.AllPods)
	info.ComputeObjectiveValues()
	return info
}

// Mutate individual via 4 operations
func (g Genetic) mutate(info *models.Individual) {
	num_pods := len(g.AllPods)
	num_nodes := len(g.AllNodes)

	// Swap two random pods assign
	swap := func(indi *models.Individual) {
		p1 := g.AllPods[rand.Intn(num_pods)]
		p2 := g.AllPods[rand.Intn(num_pods)]

		nid1 := indi.Assignment[p1.PodID]
		nid2 := indi.Assignment[p2.PodID]
		indi.Assignment[p1.PodID] = nid2
		indi.Assignment[p2.PodID] = nid1
	}

	// Change one pod assign to random one node
	change := func(indi *models.Individual) {
		pid := g.AllPods[rand.Intn(num_pods)].PodID
		nid := g.AllNodes[rand.Intn(num_nodes)].ID
		indi.Assignment[pid] = nid
	}

	// Assign one unssigned pod to a rando nodes
	assignUnassigned := func(indi *models.Individual) {
		var unassignedId []string
		for pid := range indi.Assignment {
			if indi.Assignment[pid] == "" {
				unassignedId = append(unassignedId, pid)
			}
		}
		pid := unassignedId[rand.Intn(len(unassignedId))]
		nid := g.AllNodes[rand.Intn(num_nodes)].ID
		indi.Assignment[pid] = nid
	}

	// Unset one unassigned pod
	unassignAssigned := func(indi *models.Individual) {
		var ids []string
		for pid, nid := range indi.Assignment {
			if nid != "" {
				ids = append(ids, pid)
			}
		}
		pid := ids[rand.Intn(len(ids))]
		indi.Assignment[pid] = ""
	}

	p := rand.Float64()
	if p <= 0.25 {
		change(info)
	} else if p > 0.25 && p <= 0.5 {
		swap(info)
	} else if p > 0.5 && p <= 0.51 {
		if !info.IsFeasible {
			unassignAssigned(info)
		}
	} else if p > 0.51 && p <= 1.0 {
		if info.NumOfUnassignedPods > 0 {
			assignUnassigned(info)
		}
	}
	info.ComputeObjectiveValues()
}

//******************** utils ******************************

func shuffleNodes(nodes []models.Node) {
	rand.Seed(time.Now().UnixNano())
	for n := len(nodes); n > 0; n-- {
		index := rand.Intn(n)
		nodes[n-1], nodes[index] = nodes[index], nodes[n-1]
	}
}
