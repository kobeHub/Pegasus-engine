// +build jsoniter

package main

import (
	"fmt"

	_ "github.com/kobeHub/Pegasus-engine/pkg/genetic"
	"github.com/kobeHub/Pegasus-engine/pkg/genetic/models"
)

func main() {
	fmt.Println("Welcome to pegasus-engine:", "Inno")
	var test models.Population = make([]*models.Individual, 2)
	fmt.Println(test)
}
