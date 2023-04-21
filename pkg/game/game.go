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
	MinDamage                    = 5
	MaxDamage                    = 100
)

// Color constants
const (
	clrReset       string = "\033[0m"
	clrName        string = "\033[40;1m\033[38;5;15m"
	clrHealth      string = "\033[34;1m\033[1m"
	clrDamage      string = "\033[35;1m\033[1m"
	clrGoodMessage string = "\033[32m"
	clrBadMessage  string = "\033[31m"
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

		fmt.Printf("\n%sTurn %d: %s attacks %s!%s\n\n", clrGoodMessage, currentTurn, attacker.Name, defender.Name, clrReset)
		selectedAttack := attacker.Attacks[rand.Intn(3)]
		fmt.Printf("Selected attack: %s%s%s\n", clrName, selectedAttack.Name, clrReset)

		// Determine the skill of the attacked
		attackComplexity := clamp((selectedAttack.Complexity+attacker.ComplexityBonus)/3, MinComplexity, MaxComplexity)
		fmt.Printf("Current Complexity: %.1f%%\n", attackComplexity)
		if 100*rand.Float32() > attackComplexity {
			fmt.Println("Attack performed flawlessly!")
			// Determine the attack hit chance
			attackHitChance := clamp(selectedAttack.HitChance+attacker.HitChanceBonus, MinHitChance, MaxHitChance)
			fmt.Printf("Current Hit Chance: %.1f%%\n", attackHitChance)
			if 100*rand.Float32() < attackHitChance {
				fmt.Println("Successfull hit!")
				attackBlockChance := clamp(selectedAttack.BlockChance+defender.BlockChanceBonus, MinBlockChance, MaxBlockChance)
				fmt.Printf("Current Block Chance: %.1f%%\n", attackBlockChance)
				if 100*rand.Float32() > attackBlockChance {
					fmt.Println("Attack not blocked!")
					attackDamage := clamp(selectedAttack.Damage+attacker.DamageBonus, MinDamage, MaxDamage)
					attackCriticalChance := clamp(selectedAttack.CriticalChance+attacker.CriticalChanceBonus, MinCriticalHitChance, MaxCriticalHitChance)
					fmt.Printf("Current Critical Chance: %.1f%%\n", attackBlockChance)
					if 100*rand.Float32() > attackCriticalChance {
						fmt.Println(clrDamage, "Critical hit!", clrReset)
						attackDamage = clamp(attackDamage*2, MinDamage, MaxDamage)
					}
					fmt.Printf("Damage dealt: %s%.1f%s\n", clrDamage, attackDamage, clrReset)
					defender.CurrentHealth -= int(attackDamage)
					fmt.Printf("%s%s takes %d damage! (%d/%d)%s\n", clrBadMessage, defender.Name, int(attackDamage), defender.CurrentHealth, defender.MaxHealth, clrReset)
				} else {
					fmt.Println("Attack blocked!")
				}
			} else {
				fmt.Println("Missed!")
			}
		} else {
			fmt.Println(clrBadMessage+attacker.Name, "failed to execute attack!"+clrReset)
		}

		// Switch to the next turn
		currentTurn++
		//fmt.Scanln()
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
