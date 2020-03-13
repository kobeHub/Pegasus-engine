package genetic

import (
	"math"
	"sort"

	"github.com/kobeHub/Pegasus-engine/pkg/genetic/models"
	log "github.com/sirupsen/logrus"
)

// The NSGA-II implements
type NSGA2 struct{}

// Perfoeming Fast none demainated sort
func (info NSGA2) PerformFastNonDominatedSort(popu models.Population) models.Fronts {
	var fronts models.Fronts

	for _, individual := range popu {
		// Record current `Individual` dominates states
		log.WithFields(log.Fields{
			"IsFeasible":                 individual.IsFeasible,
			"ConstraintedViolationValue": individual.ConstraintedViolationValue,
			"Objective":                  individual.ObjectiveValues,
			//"RemainingResource": individual.AllNodes,
		}).Debug("Every individual")

		for _, ano := range popu {
			if individual.ID != ano.ID {
				if individual.ConstraintDominate(*ano) {
					individual.IndividualsDominatedByThis = append(
						individual.IndividualsDominatedByThis,
						ano,
					)
				} else if ano.ConstraintDominate(*individual) {
					individual.NumOfIndividualsDominateThis++
				}
			}
		}

		//log.Info(individual.IndividualsDominatedByThis, individual.NumOfIndividualsDominateThis)

		// Check current individuals is best one
		if individual.NumOfIndividualsDominateThis == 0 {
			individual.Rank = 0
			if len(fronts) == 0 {
				fronts = append(fronts, &models.Front{})
			}
			*fronts[0] = append(*fronts[0], individual)
		}
	}

	// Add the other order `Individual`
	frontCnt := 0
	for len(*fronts[frontCnt]) != 0 {
		var next models.Front

		// Change every deminated individuals states
		for _, individual := range *fronts[frontCnt] {
			for _, dominated := range individual.IndividualsDominatedByThis {
				dominated.NumOfIndividualsDominateThis--
				if dominated.NumOfIndividualsDominateThis == 0 {
					dominated.Rank = frontCnt + 1
					next = append(next, dominated)
				}
			}
		}

		frontCnt++
		fronts = append(fronts, &next)
	}

	return fronts
}

// Compute individuals with same order distance
func (n NSGA2) ComputeCrowdingDistance(front models.Front) {
	if len(front) == 0 {
		return
	}

	for _, individual := range front {
		individual.CrowdingDistance = 0
	}

	for i := 0; i < len(front[0].ObjectiveValues); i++ {
		// Objective value descrease
		sort.Slice(front, func(first, second int) bool {
			return front[first].ObjectiveValues[i] > front[second].ObjectiveValues[i]
		})
		// First and last distance max
		front[0].CrowdingDistance = math.MaxFloat64
		front[len(front)-1].CrowdingDistance = math.MaxFloat64

		for j := 1; j < len(front)-1; j++ {
			front[j].CrowdingDistance = front[j].CrowdingDistance +
				((front[j-1].ObjectiveValues[i] - front[j+1].ObjectiveValues[i]) / (front[0].ObjectiveValues[i] - front[len(front)-1].ObjectiveValues[i]))
		}
	}
}

func (n NSGA2) SortFront(front models.Front) {
	sort.Slice(front, func(i, j int) bool {
		return !front[i].CrowdedCompareLess(*front[j])
	})
}
