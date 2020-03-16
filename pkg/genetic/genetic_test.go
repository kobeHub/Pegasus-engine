package genetic

import (
	"github.com/kobeHub/Pegasus-engine/pkg/genetic/models"
	"testing"
	"time"
)

func TestGeneticAlgorithm(t *testing.T) {
	start_time, _ := time.Parse("2006-01-02 15:04:05", "2020-03-01 01:01:01")
	n1 := models.NewConsistNode(
		"node1",
		models.NewResource(8, 16384),
		time.Now().UTC())
	n2 := models.NewConsistNode(
		"node2",
		models.NewResource(4, 9046),
		time.Now().UTC())
	n3 := models.NewConsistNode(
		"node3",
		models.NewResource(2, 1024),
		time.Now().UTC())
	n4 := models.NewDemandNode(
		"node4",
		models.NewResource(6, 16384),
		0.35,
		start_time)
	nodes := []models.Node{n1, n2, n3, n4}

	pods := []models.Pod{
		models.NewPod("1-128", 1, 128),
		models.NewPod("2-1024", 2, 1024),
		models.NewPod("2-512", 2, 512),
		models.NewPod("2-512", 2, 512),
		models.NewPod("1-1024-1", 1, 1024),
		models.NewPod("1-1024-2", 1, 1024),
		models.NewPod("0.5-128-1", 0.5, 128),
		models.NewPod("0.5-128-2", 0.5, 128),
		models.NewPod("0.5-128-3", 0.5, 128),
		models.NewPod("0.5-128-4", 0.5, 128),
	}
	g := Genetic{
		AllNodes:           nodes,
		AllPods:            pods,
		OriginalAssignment: make(map[string]string),
		Size:               212,
		GenerationNum:      1000,
	}
	t.Logf("Genetic: %v", g)

	resPopu := g.RunGeneticNSGA3(6)
	best := g.GetBestPriceIndividual(resPopu)
	t.Log(best)
	t.Log(best.Assignment)
}
