package game

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/your_username/console-game/pkg/fighter"
)

// Fight represents the fight match between two fighters
func Fight(playerFighter *fighter.Fighter, computerFighter *fighter.Fighter) *fighter.Fighter {
	rand.Seed(time.Now().UnixNano())

	currentTurn := 1 // keep track of whose turn it is
	var attacker *fighter.Fighter
	var defender *fighter.Fighter

	fmt.Printf("\n%s vs %s!\n", playerFighter.Name, computerFighter.Name)

	// Fight until one of the fighters' health is reduced to zero
	for playerFighter.CurrentHealth > 0 && computerFighter.CurrentHealth > 0 {
		// Determine who is attacking and who is defending based on the current turn
		if currentTurn%2 != 0 {
			attacker = playerFighter
			defender = computerFighter
		} else {
			attacker = computerFighter
			defender = playerFighter
		}

		fmt.Printf("\nTurn %d: %s attacks %s!\n", currentTurn, attacker.Name, defender.Name)

		// Determine the attack damage based on the attacker's speed and a random factor
		attackDamage := rand.Intn(attacker.Speed) + 1

		// Perform the attack and deduct the damage from the defender's health
		defender.CurrentHealth -= attackDamage
		fmt.Printf("%s takes %d damage! (%d/%d)\n", defender.Name, attackDamage, defender.CurrentHealth, defender.MaxHealth)

		// Switch to the next turn
		currentTurn++
	}

	// Determine the winner and return the fighter object
	var winner *fighter.Fighter
	if playerFighter.CurrentHealth > 0 {
		winner = playerFighter
	} else {
		winner = computerFighter
	}

	return winner
}