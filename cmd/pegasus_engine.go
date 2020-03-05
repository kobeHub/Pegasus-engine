// +build jsoniter

package main

import (
	"fmt"

	"github.com/kobeHub/Pegasus-engine/pkg/genetic/models"
	"github.com/kobeHub/Pegasus-engine/pkg/genetic"
)

func main() {
	fmt.Println("Welcome to pegasus-engine:", "Inno")
	var test genetic.Population = make([]*models.Individual, 2)
	fmt.Println(test)
}
