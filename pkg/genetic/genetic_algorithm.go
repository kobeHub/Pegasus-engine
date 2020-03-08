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
func (g Genetic) GenerateRandomFeasibleIndividual() *Individual {
	nodes := make([]models.Node, len(g.AllNodes))
	copy(nodes, g.AllNodes)

	// Reset Nodes remaining resources
	for i, node := range nodes {
		nodes[i].RemainingResource = node.AvailableResource
	}

	shuffleNodes(nodes)

	uid := xid.New()
	assign := make(map[string]string, len(g.AllPods))
	for _, pod := range g.AllPods {
		for _, node := range nodes {
			if pod.RequiredResource.Less(*node.RemainingResource) {
				assign[pod.PodID] = node.ID
				node.AddPod(&pod)
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
	popu := make([]*Individual, g.Size)
	for i := 0; i < g.Size; i++ {
		popu[i] = g.GenerateRandomFeasibleIndividual()
	}
	return popu
}

//******** Individual operations to generate new *********

// Select single point to crossover
func (g Genetic) reproduce(first, second Individual) Individual {
	num_pods := len(first.AllPods)
	newAssign := make(map[string]string, num_pods)
	randomCut := random.Intn(num_pods)


}


//******************** utils ******************************

func shuffleNodes(nodes []*models.Node) {
	rand.Seed(time.Now().UnixNano())
	for n := len(nodes); n > 0; n-- {
		index := rand.Intn(n)
		nodes[n-1], nodes[index] = nodes[index], nodes[n-1]
	}
}
