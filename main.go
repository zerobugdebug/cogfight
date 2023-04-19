package main

import (
	"github.com/sirupsen/logrus"

	"github.com/zerobugdebug/cogfight/pkg/fighter"
	"github.com/zerobugdebug/cogfight/pkg/game"
)

var log *logrus.Logger

func main() {
	log := logrus.New()
	// Set log format to include timestamp and colors
	log.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05.000",
		ForceColors:     true,
		FullTimestamp:   true,
	})

	// Welcome message
	log.Info("Welcome to the CogFight!")

	// Fighter Generation
	log.Info("Let's create your fighter:")
	playerFighter := fighter.CreateFighter()

	// Fight Match
	log.Info("Let's start the fight!")
	computerFighter := fighter.GenerateComputerFighter(playerFighter)
	winner := game.Fight(playerFighter, computerFighter)

	// Display the winner
	log.Info("The winner is %s!\n", winner.Name)
}
