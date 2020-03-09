package genetic

import (
	"fmt"
	"github.com/rs/xid"
	"math"
	"math/rand"
	"time"

	"github.com/kobeHub/Pegasus-engine/pkg/genetic/models"
)

// The `Population` of one generation, includes of individuals
type Population []*models.Individual

// Ordered individuals
type Front []*models.Individual
type Fronts []*Front

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
		nodes[i] = models.Node {
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
		ID: uid.String(),
		Assignment: assign,
		OriginalAssignment: g.OriginalAssignment,
	}

	info.Init(g.AllNodes, g.AllPods)
	return &info
}

func (g Genetic) GenerateRandomFeasiblePopulation() Population {
	popu := make([]*models.Individual, g.Size)
	for i := 0; i < g.Size; i++ {
		popu[i] = g.GenerateRandomFeasibleIndividual()
	}
	return popu
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
	info := models.Individual {
		ID: uid.String(),
		Assignment: newAssign,
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
	change := func(indi *INdividual) {
		pid :=
	}
}

//******************** utils ******************************

func shuffleNodes(nodes []models.Node) {
	rand.Seed(time.Now().UnixNano())
	for n := len(nodes); n > 0; n-- {
		index := rand.Intn(n)
		nodes[n-1], nodes[index] = nodes[index], nodes[n-1]
	}
}
