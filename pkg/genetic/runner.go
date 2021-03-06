package genetic

import (
	_ "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"math"

	"github.com/kobeHub/Pegasus-engine/pkg/common/k8s"
	//"github.com/kobeHub/Pegasus-engine/pkg/genetic/models"
)

func RunGeneticAlgorithm() (map[string]string, float64, error) {
	var (
		schedule   map[string]string
		totalCosts float64
	)
	pods, err := k8s.ListRescheduleLabelPods("pegasus.state/reschedulable=true")
	if err != nil {
		return schedule, totalCosts, err
	}
	nodes, err := k8s.ListReschedulableNodes(pods)
	if err != nil {
		return schedule, totalCosts, err
	}

	originalAssign := make(map[string]string, len(pods))
	for _, pod := range pods {
		originalAssign[pod.PodID] = pod.NodeID
	}

	g := Genetic{
		AllNodes:           nodes,
		AllPods:            pods,
		OriginalAssignment: originalAssign,
		Size:               viper.GetInt("PopulationSize"),
		GenerationNum:      viper.GetInt("NumOfGeneration"),
		BestPrice:          math.MaxFloat64,
	}
	_ = g.RunGeneticNSGA3(viper.GetInt("NumOfSegament"))
	best := g.GetBestPriceIndividual()
	schedule = best.Assignment
	totalCosts = best.ObjectiveValues[0]
	return schedule, totalCosts, nil
}
