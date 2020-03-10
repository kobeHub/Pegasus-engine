package genetic

import (
	"fmt"
	"github.com/rs/xid"
	"strconv"

	"github.com/kobeHub/Pegasus-engine/pkg/genetic/models"
	"github.com/kobeHub/Pegasus-engine/pkg/genetic/utils"
)

type NSGA3 struct {
	Ops NSGA2
}

func (info NSGA3) GenerateNextPopulation(t int, g Genetic, parent models.Population,
	rps []*models.ReferencePoint) models.Population {
	var (
		nextPopu models.Population
		tmpPopu  models.Population
	)
	cnt := 0
	newPopu := g.makeNewPopulation(parent)
	unionPaNew := g.combinePopulation(parent, newPopu)

	// Fast non deminate sort
	fronts := info.Ops.PerformFastNonDominatedSort(unionPaNew)
	for ;len(tmpPopu) <= g.Size; {
		info.unionFrontWithTmp(&tmpPopu, *fronts[cnt])
		cnt ++
	}

	lastFront := fronts[cnt-1]
	if len(tmpPopu) == g.Size {
		return tmpPopu
	} else {
		nextPopu = info.unionFrontsUntilLevel(fronts, cnt-1)
		num_remaining_indi := g.Size - len(nextPopu)
		utils.Normalize(tmpPopu)
		utils.Associate(tmpPopu, rps)
		info.computeNicheCount(nextPopu, rps)
		utils.Niching(num_remaining_indi, &tmpPopu, rps, lastFront, &nextPopu)
		return nextPopu
	}
}

func (info NSGA3) unionFrontWithTmp(tmp *models.Population, front models.Front) {
	*tmp = append(*tmp, front...)
}

func (info NSGA3) unionFrontsUntilLevel(fronts models.Fronts, level int) models.Population {
	var res models.Population
	for i := 0; i < level; i++ {
		res = append(res, *fronts[i]...)
	}
	return res
}

func (info NSGA3) computeNicheCount(tmpPopu models.Population, rps []*models.ReferencePoint) {
	for _, indi := range tmpPopu {
		for _, rp := range rps {
			if indi.ReferencePoint.ID == rp.ID {
				rp.NicheCount ++
			}
		}
	}
}

// reference point coordinates
var rpsCoordinates map[string]models.ReferencePoint

func (info NSGA3) GetReferencePoints(num_objective, num_segament int) []*models.ReferencePoint {
	rpsc := map[string]models.ReferencePoint{}
	init_rpsc := make([]float64, num_objective)
	init_rpsc[0] = 1.
	init_rps := models.ReferencePoint{
		Coordinates: init_rpsc,
	}
	info.generateRPSC(init_rps, num_objective, num_segament)

	var result []*models.ReferencePoint
	for _, item := range rpsc {
		uid := xid.New()
		result = append(result, &models.ReferencePoint{
			ID: uid.String(),
			Coordinates: item.Coordinates,
		})
	}
	return result
}


func (info NSGA3) generateRPSC(rp models.ReferencePoint, num_objective, num_segament int) int {
	if rp.Coordinates[0] < 0 {
		return 0
	}
	rpsCoordinates[fmt.Sprint(rp.Coordinates)] = rp

	for i := 1; i < num_objective; i++ {
		newCoordinates := make([]float64, len(rp.Coordinates))
		copy(newCoordinates, rp.Coordinates)
		newCoordinates[0] = newCoordinates[0] - info.Round(1. / float64(num_segament))
		newCoordinates[i] = newCoordinates[i] + info.Round(1. / float64(num_segament))
		if _, exist := rpsCoordinates[fmt.Sprint(newCoordinates)]; !exist {
			newRP := models.ReferencePoint{Coordinates: newCoordinates}
			info.generateRPSC(newRP, num_objective, num_segament)
		}
	}

	return 0
}

func (info NSGA3) Round(num float64) float64 {
	float, _:= strconv.ParseFloat(fmt.Sprint("%.2f", num), 64)
	return float
}
