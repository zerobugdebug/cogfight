package game

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/zerobugdebug/cogfight/pkg/fighter"
)

const (
	MaxHitChance         float32 = 95
	MinHitChance                 = 5
	MaxBlockChance               = 95
	MinBlockChance               = 5
	MaxComplexity                = 95
	MinComplexity                = 5
	MaxCriticalHitChance         = 95
	MinCriticalHitChance         = 5
	MaxSpecialChance             = 95
	MinSpecialChance             = 5
)

func clamp(val, min, max float32) float32 {
	if val < min {
		return min
	} else if val > max {
		return max
	} else {
		return val
	}
}

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

		fmt.Printf("\nTurn %d: %s attacks %s!\n\n", currentTurn, attacker.Name, defender.Name)
		selectedAttack := attacker.Attacks[rand.Intn(3)]
		fmt.Printf("Selected attack: %s\n", selectedAttack.Name)

		// Determine the attack hit chance
		attackHitChance := clamp(selectedAttack.HitChance+attacker.HitChanceBonus, MinHitChance, MaxHitChance)
		fmt.Printf("Current Hit Chance: %.1f%%\n", attackHitChance)
		if 100*rand.Float32() < attackHitChance {
			fmt.Println("Successfull hit!")
			attackBlockChance := clamp(selectedAttack.BlockChance+defender.BlockChanceBonus, MinBlockChance, MaxBlockChance)
			fmt.Printf("Current Block Chance: %.1f%%\n", attackBlockChance)
			if 100*rand.Float32() > attackBlockChance {
				fmt.Println("Attack not blocked!")
				attackDamage := selectedAttack.Damage + attacker.DamageBonus
				fmt.Printf("Damage dealt: %.1f%%\n", attackDamage)
				defender.CurrentHealth -= int(attackDamage)
				fmt.Printf("%s takes %d damage! (%d/%d)\n", defender.Name, int(attackDamage), defender.CurrentHealth, defender.MaxHealth)
			} else {
				fmt.Println("Attack blocked!")
			}
		} else {
			fmt.Println("Missed!")
		}

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
