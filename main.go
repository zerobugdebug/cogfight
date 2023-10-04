package main

import (
	"github.com/zerobugdebug/cogfight/pkg/fighter"
	"github.com/zerobugdebug/cogfight/pkg/game"
	"github.com/zerobugdebug/cogfight/pkg/logging"

)

func main() {
	//fmt.Println("log = %v", log)
	// Welcome message
	logging.Info("Welcome to the CogFight!")

	// Fighter Generation
	logging.Info("Let's create your fighter:")
	playerFighter := fighter.CreateFighter()

	// Fight Match
	logging.Info("Let's start the fight!")
	computerFighter := fighter.GenerateComputerFighter(playerFighter)
	winner := game.Fight(playerFighter, computerFighter)
	if winner != nil {
		// Display the winner
		logging.Info("The winner is ", winner.Name)
	}
}
