package game

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"

	//"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"

	"github.com/zerobugdebug/cogfight/pkg/attack"
	"github.com/zerobugdebug/cogfight/pkg/fighter"
	"github.com/zerobugdebug/cogfight/pkg/modifiers"
	"github.com/zerobugdebug/cogfight/pkg/ui"
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
	fighter.DisplayFighters(playerFighter, computerFighter)
	fmt.Println("Waiting for the comments...")
	stopChan := make(chan bool)
	var wg sync.WaitGroup
	wg.Add(1)
	go ui.RotatingPipe(stopChan, &wg)
	situation := fmt.Sprintf("Fight not started yet. Commentators introduce themselves and talk about the fighters\nFirst fighter: %s Second fighter: %s", playerFighter.String(), computerFighter.String())
	var chatMessages []fighter.ChatMessage = []fighter.ChatMessage{{Role: "user", Content: situation}}
	comments, err := fighter.GetOpenAIResponse("COG_TURN_COMMENT_PROMPT", chatMessages, "full")
	if err != nil {
		fmt.Println("Can't get OpenAI response")
		return nil
	}
	stopChan <- true
	wg.Wait()
	strComments := strings.Replace(comments.(string), "\n\n", "\n", -1)
	fmt.Println("\n" + strComments)
	chatMessages = append(chatMessages, fighter.ChatMessage{Role: "assistant", Content: strComments})
	fmt.Scanln()

	var situationDescription string
	//situationDescription = fmt.Sprintf("First fighter: %s Second fighter: %s", playerFighter.String(), computerFighter.String())
	//var chatMessages []fighter.ChatMessage = []fighter.ChatMessage{{Role: "system", Content: situationDescription}}
	// Fight until one of the fighters' health is reduced to zero
	for playerFighter.CurrentHealth > 0 && computerFighter.CurrentHealth > 0 {
		//		prevSituationDescription = "Previous rounds: \n" + situationDescription + "\n Current round to be described: \n"
		//prevSituationDescription = situationDescription
		situationDescription = ""
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
						situationDescription += attacker.Name + " is currently " + condition.String() + ". "
					}
				}
			}
		}

		if skipTurn != 0 {
			fmt.Printf("\n%sTurn %d: %s cannot attack, skipping turn!%s\n\n", clrGoodMessage, currentTurn, attacker.Name, clrReset)
			situationDescription += attacker.Name + " cannot attack. "
		} else {
			var selectedAttack *attack.Attack
			fmt.Printf("\n%sTurn %d: %s attacks %s!%s\n\n", clrGoodMessage, currentTurn, attacker.Name, defender.Name, clrReset)
			if currentTurn%2 != 0 {
				selectedAttack = attacker.SelectAttack(defender)
			} else {
				selectedAttack = attack.NewDefaultAttacks().GetRandomAttack()
				fmt.Printf("Selected attack: %s\n", color.CyanString(selectedAttack.Name))
			}
			//situationDescription += attacker.Name + " executing " + selectedAttack.Name + ". "
			attackResultText := attacker.ApplyAttack(defender, selectedAttack)
			situationDescription += attackResultText
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
				situationDescription += attacker.Name + " is not " + condition.String() + " anymore. "
			}

		}
		if defender.CurrentHealth < 0 {
			situationDescription += defender.Name + " is knocked out. "
		}
		if attacker.CurrentHealth < 0 {
			situationDescription += attacker.Name + " lost consciousness. "
		}
		wg.Add(1)
		go ui.RotatingPipe(stopChan, &wg)
		//situation := fmt.Sprintf("Previous rounds:\n %s\n Current round to be described:\nTurn %d: %s attacks %s. %s", prevSituationDescription, currentTurn, attacker.Name, defender.Name, situationDescription)
		situation := fmt.Sprintf("Turn %d: %s attacks %s. %s", currentTurn, attacker.Name, defender.Name, situationDescription)
		chatMessages = append(chatMessages, fighter.ChatMessage{Role: "user", Content: situation})
		// fmt.Printf("\n-------------------------------\n")
		// fmt.Printf("prevSituationDescription: %v\n", prevSituationDescription)
		// fmt.Printf("\n-------------------------------\n")
		// fmt.Printf("situationDescription: %v\n", situationDescription)
		// fmt.Printf("\n-------------------------------\n")
		// fmt.Printf("situation: %v\n", situation)
		comments, err := fighter.GetOpenAIResponse("COG_TURN_COMMENT_PROMPT", chatMessages, "full")
		if err != nil {
			fmt.Println("Can't get OpenAI response")
			return nil
		}
		stopChan <- true
		wg.Wait()
		strComments := strings.Replace(comments.(string), "\n\n", "\n", -1)
		fmt.Println("\n" + strComments)
		chatMessages = append(chatMessages, fighter.ChatMessage{Role: "assistant", Content: strComments})
		//prevSituationDescription = fmt.Sprintf("%s\nTurn %d: %s attacks %s. \n%s\n%s\n", prevSituationDescription, currentTurn, attacker.Name, defender.Name, situationDescription, comments.(string))
		//fmt.Printf("\n-------------------------------\n")
		//fmt.Printf("chatMessages: %v\n", chatMessages)
		//fmt.Printf("\n-------------------------------\n")
		currentTurn++
		fmt.Scanln()
	}

	// Determine the winner and return the fighter object
	var winner *fighter.Fighter
	if computerFighter.CurrentHealth < 0 {
		winner = playerFighter
	} else {
		winner = computerFighter
	}

	return winner
}
