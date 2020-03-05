package genetic

import (
	"github.com/kobeHub/Pegasus-engine/pkg/genetic/models"
)

// The `Population` of one generation, includes of individuals
type Population []*models.Individual

// Ordered individuals
type Front []*models.Individual
type Fronts []*Front
