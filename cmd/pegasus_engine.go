// +build jsoniter

package main

import (
	"fmt"

	"github.com/kobeHub/Pegasus-engine/pkg/genetic"
	"github.com/kobeHub/Pegasus-engine/pkg/genetic/models"
	_ "github.com/kobeHub/Pegasus-engine/pkg/genetic/utils"
)

func main() {
	fmt.Println("Welcome to pegasus-engine:", "Inno")
	var test genetic.Population = make([]*models.Individual, 2)
	fmt.Println(test)
}
