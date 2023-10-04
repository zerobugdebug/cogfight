package main

import (
	"github.com/zerobugdebug/cogfight/pkg/fighter"
	"github.com/zerobugdebug/cogfight/pkg/game"
	"github.com/zerobugdebug/cogfight/pkg/log"
)

func main() {

	// Welcome message
	log.Info("Welcome to the CogFight!")

	// Fighter Generation
	log.Info("Let's create your fighter:")
	playerFighter := fighter.CreateFighter()

	// Fight Match
	log.Info("Let's start the fight!")
	computerFighter := fighter.GenerateComputerFighter(playerFighter)
	winner := game.Fight(playerFighter, computerFighter)
	if winner != nil {
		// Display the winner
		log.Info("The winner is ", winner.Name)
	}
}
