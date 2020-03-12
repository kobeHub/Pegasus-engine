package genetic

import (
	"github.com/kobeHub/Pegasus-engine/pkg/genetic/models"
	"github.com/rs/xid"
	"github.com/spf13/viper"

	"testing"
	"time"
)

func TestGeneticAlgorithm(t *testing.T) {
	start_time, _ := time.Parse("2006-0102 15:04:05", "2020-03-01 01:01:01")
	n1 := models.NewConsistNode(
		xid.New().String(),
		models.NewResource(8, 16384))
	n2 := models.NewConsistNode(
		xid.New().String(),
		models.NewResource(4, 9046))
	n3 := models.NewConsistNode(
		xid.New().String(),
		models.NewResource(2, 1024))
	n4 := models.NewDemandNode(
		xid.New().String(),
		models.NewResource(6, 16384),
		0.35,
		start_time)
	nodes := []models.Node{n1, n2, n3, n4}

	pods := []models.Pod{
		models.NewPod(xid.New().String(), 1, 128),
		models.NewPod(xid.New().String(), 2, 1024),
		models.NewPod(xid.New().String(), 2, 512),
		models.NewPod(xid.New().String(), 2, 512),
		models.NewPod(xid.New().String(), 1, 1024),
		models.NewPod(xid.New().String(), 1, 1024),
		models.NewPod(xid.New().String(), 0.5, 128),
		models.NewPod(xid.New().String(), 0.5, 128),
		models.NewPod(xid.New().String(), 0.5, 128),
		models.NewPod(xid.New().String(), 0.5, 128),
	}
	g := Genetic{
		AllNodes: nodes,
		AllPods: pods,
		OriginalAssignment: make(map[string]string),
		Size: 10,
		GenerationNum: 6,
	}
	// t.Logf("Genetic: %v", g)

	resPopu := g.RunGeneticNSGA3(viper.GetInt("NumOfSegament"))
	t.Log(resPopu)
}
