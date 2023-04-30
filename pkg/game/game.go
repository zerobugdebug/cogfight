package game

import (
	"fmt"
	"math/rand"
	"time"

	//"github.com/AlecAivazis/survey/v2"

	"github.com/zerobugdebug/cogfight/pkg/attack"
	"github.com/zerobugdebug/cogfight/pkg/fighter"
	"github.com/zerobugdebug/cogfight/pkg/modifiers"
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
		fighter.DisplayFighters(playerFighter, computerFighter)
		skipTurn := 0

		//Apply pre-turn conditions
		for condition := range attacker.Conditions {
			for modifier, value := range modifiers.DefaultConditionAttributes[condition] {
				switch modifier {
				case modifiers.SkipTurn:
					{
						skipTurn = value
					}
				}
			}
		}

		if skipTurn != 0 {
			fmt.Printf("\n%sTurn %d: %s cannot attack, skipping turn!%s\n\n", clrGoodMessage, currentTurn, attacker.Name, clrReset)
		} else {
			var selectedAttack *attack.Attack
			fmt.Printf("\n%sTurn %d: %s attacks %s!%s\n\n", clrGoodMessage, currentTurn, attacker.Name, defender.Name, clrReset)
			if currentTurn%2 != 0 {
				selectedAttack = attacker.SelectAttack(defender)
			} else {
				selectedAttack = attack.NewDefaultAttacks().GetRandomAttack()
			}
			//fmt.Printf("Selected attack: %s\n", color.CyanString(selectedAttack.Name))
			attacker.ApplyAttack(defender, selectedAttack)
		}

		//Apply post-turn conditions
		//Calculate effect from attacker conditions
		for condition := range attacker.Conditions {
			for modifier, value := range modifiers.DefaultConditionAttributes[condition] {
				switch modifier {
				case modifiers.HPPerTurn:
					{
						attacker.CurrentHealth += int(value)
						if int(value) < 0 {
							fmt.Printf("%s takes %d damage! (%d/%d) due to %s\n", attacker.Name, -int(value), attacker.CurrentHealth, attacker.MaxHealth, condition.String())
						}
					}
				}
			}
			attacker.Conditions[condition] -= 1
			if attacker.Conditions[condition] < 1 {
				delete(attacker.Conditions, condition)
				defender.RemoveCondition(attacker, condition)
			}

		}

		currentTurn++
		fmt.Scanln()
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
