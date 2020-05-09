package genetic

import (
	"fmt"
	"github.com/kobeHub/Pegasus-engine/pkg/genetic/models"
	"math"
	"testing"
	"time"
)

func allNodes() []models.Node {
	nodes := make([]models.Node, 0, 24)
	for i := 0; i < 6; i++ {
		id := fmt.Sprintf("g6.large.%v", i)
		nodes = append(nodes, models.NewConsistNode(
			id,
			id,
			models.NewResource(2, 8192),
			time.Now().UTC()))
	}
	for i := 0; i < 4; i++ {
		id := fmt.Sprintf("g6.xlarge.%v", i)
		nodes = append(nodes, models.NewDemandNode(
			id,
			id,
			models.NewResource(4, 16384),
			0.7,
			time.Now().UTC().Add(time.Hour * (-24))))
	}
	for i := 0; i < 5; i++ {
		id := fmt.Sprintf("g6.2xlarge.%v", i)
		nodes = append(nodes, models.NewConsistNode(
			id,
			id,
			models.NewResource(8, 32768),
			time.Now().UTC()))
	}
	for i := 0; i < 3; i++ {
		id := fmt.Sprintf("r6.large.%v", i)
		nodes = append(nodes, models.NewDemandNode(
			id,
			id,
			models.NewResource(2, 16384),
			0.46,
			time.Now().UTC().Add(time.Hour * (-24))))
	}
	for i := 0; i < 5; i++ {
		id := fmt.Sprintf("r6.xlarge.%v", i)
		nodes = append(nodes, models.NewConsistNode(
			id,
			id,
			models.NewResource(4, 32768),
			time.Now().UTC()))
	}
	nodes = append(nodes,  models.NewDemandNode(
		"r6.2xlarge.1",
		"r6.2xlarge.1",
		models.NewResource(8, 65536),
		1.83,
		time.Now().UTC().Add(time.Hour * (-24))))
	return nodes
}

func generatePods() []models.Pod {
	pods := make([]models.Pod, 0, 100)
	for i := 0; i < 3; i++ {
		pods = append(pods, models.NewPod(
			fmt.Sprintf("6-12-%v", i),
			6,
			12 * 1024,
			"",
			models.Unknown,
		))
		pods = append(pods, models.NewPod(
			fmt.Sprintf("3-6-%v", i),
			3,
			6 * 1024,
			"",
			models.Unknown,
		))
	}
	for i := 0; i < 10; i++ {
		pods = append(pods, models.NewPod(
			fmt.Sprintf("1-2-%v", i),
			1,
			2 * 1024,
			"",
			models.Unknown,
		))
	}
	for i := 0; i < 5; i++ {
		pods = append(pods, models.NewPod(
			fmt.Sprintf("2-4-%v", i),
			2,
			4 * 1024,
			"",
			models.Unknown,
		))
	}
	for i := 0; i < 50; i++ {
		pods = append(pods, models.NewPod(
			fmt.Sprintf("0.5-1-%v", i),
			0.5,
			1024,
			"",
			models.Unknown,
		))
	}
	return pods
}

func TestGeneticAlgorithm(t *testing.T) {
	nodes := allNodes()
	pods := generatePods()
	t.Log("Nodes: ", len(nodes))
	t.Log("Pods: ", len(pods))

	size :=100
	num := 600
	g := Genetic{
		AllNodes:           nodes,
		AllPods:            pods,
		OriginalAssignment: make(map[string]string),
		Size:               size,
		GenerationNum:      num,
		BestPrice:          math.MaxFloat64,
	}
	t.Logf("Genetic: %v", g)


	_ = g.RunGeneticNSGA3(6)
	best := g.GetBestPriceIndividual()
	t.Logf("Size: %v, Generation: %v\n", size, num)
	t.Log("Price:", best.ObjectiveValues[0])
	t.Log(best.Assignment)
}
