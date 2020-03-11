// +build jsoniter

package main

import (
	"os"

	"github.com/kobeHub/Pegasus-engine/config"
	_ "github.com/kobeHub/Pegasus-engine/pkg/genetic"
	"github.com/kobeHub/Pegasus-engine/pkg/genetic/models"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	if len(os.Args) != 2 {
		panic("Please spefic environment:\nlocal, test or prod")
	}
	config.Init(os.Args[1])
	log.Info("Welcome to pegasus-engine")
	var test models.Population = make([]*models.Individual, 2)
	log.Info(test)
	log.Info(viper.GetInt("PopulationSize"))
}
