package main

import (
	"fmt"

	"github.com/zerobugdebug/cogfight/pkg/fighter"
	"github.com/zerobugdebug/cogfight/pkg/game"
)

func main() {
	// Welcome message
	fmt.Println("Welcome to the Console Game!")

	// Fighter Generation
	fmt.Println("\nLet's create your fighter:")
	playerFighter := fighter.CreateFighter()

	// Fight Match
	fmt.Println("\nLet's start the fight!")
	computerFighter := fighter.GenerateComputerFighter(playerFighter)
	winner := game.Fight(playerFighter, computerFighter)

	// Display the winner
	fmt.Printf("\nThe winner is %s!\n", winner.Name)
}