package fighter

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"

	"github.com/AlecAivazis/survey/v2"
	"github.com/sashabaranov/go-openai"
)

// Attack represents an attack in the game
type Attack struct {
	Name           string
	Damage         int
	Complexity     int
	HitChance      int
	BlockChance    int
	CriticalChance int
	TakedownChance int
}

// Fighter represents a fighter in the game
type Fighter struct {
	Name          string
	Height        int
	Weight        int
	Age           int
	Speed         int
	Attacks       []*Attack
	CurrentHealth int
	MaxHealth     int
}

// CreateFighter creates a new fighter object based on user input
func CreateFighter() *Fighter {
	// Collect user input
	var name string
	namePrompt := &survey.Input{
		Message: "Enter fighter name:",
	}
	survey.AskOne(namePrompt, &name)

	var height int
	heightPrompt := &survey.Input{
		Message: "Enter fighter height:",
	}
	survey.AskOne(heightPrompt, &height)

	var weight int
	weightPrompt := &survey.Input{
		Message: "Enter fighter weight:",
	}
	survey.AskOne(weightPrompt, &weight)

	var age int
	agePrompt := &survey.Input{
		Message: "Enter fighter age:",
	}
	survey.AskOne(agePrompt, &age)

	var speed int
	speedPrompt := &survey.Input{
		Message: "Enter fighter speed:",
	}
	survey.AskOne(speedPrompt, &speed)

	attacks := []*Attack{}
	attackPrompt := &survey.Input{
		Message: "Enter an attack name (leave empty to finish):",
	}
	for {
		var attackName string
		err := survey.AskOne(attackPrompt, &attackName, survey.WithValidator(survey.Required))
		if err != nil {
			break
		}

		// Validate the attack name and get the attack parameters using OpenAI API
		attack, err := validateAttackName(attackName)
		if err != nil {
			fmt.Printf("Error validating attack name: %s\n", err)
			continue
		}

		attacks = append(attacks, attack)
	}

	// Create the fighter object
	fighter := &Fighter{
		Name:          name,
		Height:        height,
		Weight:        weight,
		Age:           age,
		Speed:         speed,
		Attacks:       attacks,
		CurrentHealth: 100,
		MaxHealth:     100,
	}

	fmt.Printf("\n%s has been created!\n", fighter.Name)
	return fighter
}

// GenerateComputerFighter generates a computer-controlled fighter
func GenerateComputerFighter(playerFighter *Fighter) *Fighter {
	// Generate random values for the computer fighter's attributes
	computerFighter := &Fighter{
		Name:          "Computer Fighter",
		Height:        rand.Intn(50) + 150, // Height between 150 and 199 cm
		Weight:        rand.Intn(50) + 50,  // Weight between 50 and 99 kg
		Age:           rand.Intn(30) + 20,  // Age between 20 and 49 years
		Speed:         rand.Intn(10) + 1,   // Speed between 1 and 10 m/s
		Attacks:       []*Attack{},
		CurrentHealth: 100,
		MaxHealth:     100,
	}

	// Copy the player's attacks and modify the parameters for the computer's attacks
	for _, playerAttack := range playerFighter.Attacks {
		computerAttack := &Attack{
			Name:           playerAttack.Name,
			Damage:         playerAttack.Damage - 10,
			Complexity:     playerAttack.Complexity + 1,
			HitChance:      playerAttack.HitChance - 5,
			BlockChance:    playerAttack.BlockChance - 5,
			CriticalChance: playerAttack.CriticalChance - 5,
			TakedownChance: playerAttack.TakedownChance - 5,
		}

		computerFighter.Attacks = append(computerFighter.Attacks, computerAttack)
	}

	fmt.Printf("\n%s has been generated!\n", computerFighter.Name)
	return computerFighter
}

// validateAttackName validates the given attack name using OpenAI API and returns the attack parameters
func validateAttackName(attackName string) (*Attack, error) {
	client := openai.NewClient("YOUR_API_KEY_HERE")

	// Define the prompt template
	promptTemplate := `{
		"prompt": "Please provide the parameters for the '%s' attack:",
		"max_tokens": 100,
		"model": "text-davinci-002",
		"temperature": 0.5,
		"stop": ["\n"]
	}`

	// Send the prompt to OpenAI API and get the response
	prompt := fmt.Sprintf(promptTemplate, attackName)
	params := openai.CompletionRequest{
		Prompt: &prompt,
	}
	response, err := client.CreateCompletion(context.Background(), params)
	if err != nil {
		return nil, fmt.Errorf("error sending OpenAI API request: %s", err)
	}

	// Parse the response to get the attack parameters
	paramsJSON := response.Choices[0].Text
	attack := &Attack{}
	err = json.Unmarshal([]byte(paramsJSON), attack)
	if err != nil {
		return nil, fmt.Errorf("error parsing OpenAI API response: %s", err)
	}

	attack.Name = attackName
	return attack, nil
}
