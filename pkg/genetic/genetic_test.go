package genetic

import (
	"github.com/kobeHub/Pegasus-engine/pkg/genetic/models"
	"math"
	"testing"
	"time"
)

func TestGeneticAlgorithm(t *testing.T) {
	n1 := models.NewConsistNode(
		"node1",
		"node1",
		models.NewResource(2, 8192),
		time.Now().UTC())
	n2 := models.NewConsistNode(
		"node2",
		"node2",
		models.NewResource(4, 16384),
		time.Now().UTC())
	n3 := models.NewConsistNode(
		"node3",
		"node3",
		models.NewResource(8, 32768),
		time.Now().UTC())
	n4 := models.NewConsistNode(
		"node4",
		"node4",
		models.NewResource(12, 49152),
		time.Now().UTC())
	n5 := models.NewDemandNode(
		"node5",
		"node5",
		models.NewResource(8, 32768),
		1.4,
		time.Now().UTC().Add(time.Hour * (-24)))
	n6 := models.NewDemandNode(
		"node6",
		"node6",
		models.NewResource(4, 16384),
		0.7,
		time.Now().UTC().Add(time.Hour * (-24)))
	n7 := models.NewDemandNode(
		"node7",
		"node7",
		models.NewResource(2, 8192),
		0.35,
		time.Now().UTC().Add(time.Hour * (-24)))
	nodes := []models.Node{n1, n2, n3, n4, n5, n6, n7}

	pods := []models.Pod{
		models.NewPod("1-500", 1, 500, "", models.Unknown),
		models.NewPod("2-2048", 2, 2048, "", models.Unknown),
		models.NewPod("2-4096", 2, 4096, "", models.Unknown),
		models.NewPod("6-10240", 6, 10240, "", models.Unknown),
		models.NewPod("1-1024-1", 1, 1024, "", models.Unknown),
		models.NewPod("1-1024-2", 1, 1024, "", models.Unknown),
		models.NewPod("0.5-512-1", 0.5, 512, "", models.Unknown),
		models.NewPod("0.5-512-2", 0.5, 512, "", models.Unknown),
		models.NewPod("0.5-512-3", 0.5, 512, "", models.Unknown),
		models.NewPod("0.5-512-4", 0.5, 512, "", models.Unknown),
		models.NewPod("1-2048", 1, 2048, "", models.Unknown),
		models.NewPod("3-5120", 3, 5120, "", models.Unknown),
	}
	g := Genetic{
		AllNodes:           nodes,
		AllPods:            pods,
		OriginalAssignment: make(map[string]string),
		Size:               400,
		GenerationNum:      1000,
		BestPrice:          math.MaxFloat64,
	}
	t.Logf("Genetic: %v", g)

	_ = g.RunGeneticNSGA3(6)
	best := g.GetBestPriceIndividual()
	t.Log("Price:", best.ObjectiveValues[0])
	t.Log(best.Assignment)
}
