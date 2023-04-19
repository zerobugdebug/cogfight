package fighter

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strconv"

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
	Name                        string
	Height                      int
	Weight                      int
	Age                         int
	AgilityStrengthBalance      float32
	BurstEnduranceBalance       float32
	DefenseOffenseBalance       float32
	SpeedControlBalance         float32
	IntelligenceInstinctBalance float32
	DamageBonus                 float32
	ComplexityBonus             float32
	HitChanceBonus              float32
	BlockChanceBonus            float32
	CriticalChanceBonus         float32
	TakedownChanceBonus         float32
	Attacks                     []*Attack
	CurrentHealth               int
	MaxHealth                   int
}

// validateNumber requires that the number provided was between min and max
func validateNumber(optParams ...int) survey.Validator {
	var min, max int
	switch len(optParams) {
	case 1:
		min = optParams[0]
		max = int(^uint(0) >> 1) // set max to the maximum value of int
	case 2:
		min = optParams[0]
		max = optParams[1]
	}
	// return a validator that checks the length of the list
	return func(val interface{}) error {
		if str, ok := val.(string); !ok {
			return errors.New("Answer should be a string")
		} else {
			height, err := strconv.Atoi(str)
			if err != nil {
				return errors.New("Answer should be a number")
			}
			if len(optParams) > 0 && (height < min || height > max) {
				return fmt.Errorf("Answer should be a number between %d and %d", min, max)
			}
		}
		return nil
	}

}

// CreateFighter creates a new fighter object based on user input
func CreateFighter() *Fighter {
	// Collect user input
	// Define the survey questions array
	qs := []*survey.Question{}

	// Define the survey questions
	nameQuestion := &survey.Question{
		Name: "name",
		Prompt: &survey.Input{
			Message: "Enter fighter name:",
		},
		Validate: survey.Required,
	}
	qs = append(qs, nameQuestion)

	heightQuestion := &survey.Question{
		Name: "height",
		Prompt: &survey.Input{
			Message: "Enter fighter height (150-220 cm):",
		},
		Validate: validateNumber(150, 220),
	}
	qs = append(qs, heightQuestion)

	weightQuestion := &survey.Question{
		Name: "weight",
		Prompt: &survey.Input{
			Message: "Enter fighter weight (50-200 kg):",
		},
		Validate: validateNumber(50, 200),
	}
	qs = append(qs, weightQuestion)

	ageQuestion := &survey.Question{
		Name: "age",
		Prompt: &survey.Input{
			Message: "Enter fighter age (18-60 years):",
		},
		Validate: validateNumber(18, 60),
	}
	qs = append(qs, ageQuestion)

	agilityStrengthBalanceQuestion := &survey.Question{
		Name: "agilityStrengthBalance",
		Prompt: &survey.Select{
			Message: "Choose fighter agility/strength balance:",
			Options: []string{"Very high Agility, Very low Strength", "High Agility, Low Strength", "Balanced", "Low Agility, High Strength", "Very low Agility, Very high Strength"},
		},
	}
	qs = append(qs, agilityStrengthBalanceQuestion)

	burstEnduranceBalanceQuestion := &survey.Question{
		Name: "burstEnduranceBalance",
		Prompt: &survey.Select{
			Message: "Choose fighter burst/endurance balance:",
			Options: []string{"Very high Burst, Very low Endurance", "High Burst, Low Endurance", "Balanced", "Low Burst, High Endurance", "Very low Burst, Very high Endurance"},
		},
	}
	qs = append(qs, burstEnduranceBalanceQuestion)

	defenseOffenseBalanceQuestion := &survey.Question{
		Name: "defenseOffenseBalance",
		Prompt: &survey.Select{
			Message: "Choose fighter defense/offense balance:",
			Options: []string{"Very high Defense, Very low Offense", "High Defense, Low Offense", "Balanced", "Low Defense, High Offense", "Very low Defense, Very high Offense"},
		},
	}
	qs = append(qs, defenseOffenseBalanceQuestion)

	speedControlBalanceQuestion := &survey.Question{
		Name: "speedControlBalance",
		Prompt: &survey.Select{
			Message: "Choose fighter speed/control balance:",
			Options: []string{"Very high Speed, Very low Control", "High Speed, Low Control", "Balanced", "Low Speed, High Control", "Very low Speed, Very high Control"},
		},
	}
	qs = append(qs, speedControlBalanceQuestion)

	intelligenceInstinctBalanceQuestion := &survey.Question{
		Name: "intelligenceInstinctBalance",
		Prompt: &survey.Select{
			Message: "Choose fighter intelligence/instinct balance:",
			Options: []string{"Very high Intelligence, Very low Instinct", "High Intelligence, Low Instinct", "Balanced", "Low Intelligence, High Instinct", "Very low Intelligence, Very high Instinct"},
		},
	}
	qs = append(qs, intelligenceInstinctBalanceQuestion)

	// Ask the user for input
	answers := struct {
		Name                        string
		Height                      int
		Weight                      int
		Age                         int
		AgilityStrengthBalance      int
		BurstEnduranceBalance       int
		DefenseOffenseBalance       int
		SpeedControlBalance         int
		IntelligenceInstinctBalance int
	}{}

	err := survey.Ask(qs, &answers)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	fmt.Println(answers)

	attacks := []*Attack{}
	attackPrompt := &survey.Input{
		Message: "Enter an attack name:",
	}

	for i := 0; i < 3; i++ {
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
		Name:                        answers.Name,
		Height:                      answers.Height,
		Weight:                      answers.Weight,
		Age:                         answers.Age,
		AgilityStrengthBalance:      float32(answers.AgilityStrengthBalance) + (float32(answers.Weight)-125)/50 + (float32(answers.Height)-185)/20,
		BurstEnduranceBalance:       float32(answers.BurstEnduranceBalance) - (float32(answers.Weight)-125)/50,
		DefenseOffenseBalance:       float32(answers.DefenseOffenseBalance) + (float32(answers.Height)-185)/20,
		SpeedControlBalance:         float32(answers.SpeedControlBalance) + (float32(answers.Weight)-125)/50 + (float32(answers.Height)-185)/20,
		IntelligenceInstinctBalance: float32(answers.IntelligenceInstinctBalance) - (float32(answers.Age)-39)/10,
		CurrentHealth:               100,
		MaxHealth:                   100,
	}

	fighter.DamageBonus = (fighter.AgilityStrengthBalance + fighter.BurstEnduranceBalance - 4) * 2
	fighter.ComplexityBonus = (-fighter.AgilityStrengthBalance + fighter.SpeedControlBalance - fighter.IntelligenceInstinctBalance + 2) * 2
	fighter.HitChanceBonus = (-fighter.AgilityStrengthBalance - fighter.BurstEnduranceBalance + fighter.DefenseOffenseBalance - fighter.SpeedControlBalance + fighter.IntelligenceInstinctBalance + 2) * 2
	fighter.BlockChanceBonus = (-fighter.AgilityStrengthBalance + fighter.BurstEnduranceBalance - fighter.DefenseOffenseBalance - fighter.SpeedControlBalance + fighter.IntelligenceInstinctBalance + 2) * 2
	fighter.CriticalChanceBonus = (fighter.SpeedControlBalance - fighter.IntelligenceInstinctBalance) * 2
	fighter.TakedownChanceBonus = (fighter.AgilityStrengthBalance - fighter.BurstEnduranceBalance) * 2

	fmt.Println(fighter)
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
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("OpenAI API key not found in environment variable OPENAI_API_KEY")
	}

	client := openai.NewClient(apiKey)

	// Define the prompt template
	promptTemplate := os.Getenv("COG_VALIDATION_ATTACK_PROMPT")

	// Send the prompt to OpenAI API and get the response
	prompt := fmt.Sprintf(promptTemplate, attackName)

	response, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:     openai.GPT3Dot5Turbo,
			MaxTokens: 3,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
		},
	)
	if err != nil {
		return nil, fmt.Errorf("error sending OpenAI API request: %s", err)
	}

	// Parse the response to confirm if attack is valid
	paramsJSON := response.Choices[0].Message.Content
	fmt.Println(paramsJSON)
	confirmed := false
	confPrompt := &survey.Confirm{
		Message: "Confirmed?",
	}
	survey.AskOne(confPrompt, &confirmed)

	attack := &Attack{}
	err = json.Unmarshal([]byte(paramsJSON), attack)
	if err != nil {
		return nil, fmt.Errorf("error parsing OpenAI API response: %s", err)
	}

	attack.Name = attackName
	return attack, nil
}

// SaveFighterToFile saves a fighter object to a JSON file
func SaveFighterToFile(fighter *Fighter, filename string) error {
	// Convert the fighter object to JSON
	fighterJSON, err := json.MarshalIndent(fighter, "", "  ")
	if err != nil {
		return fmt.Errorf("error encoding fighter to JSON: %s", err)
	}

	// Write the JSON data to the file
	err = os.WriteFile(filename, fighterJSON, 0644)
	if err != nil {
		return fmt.Errorf("error writing fighter data to file: %s", err)
	}

	fmt.Printf("Fighter data saved to %s!\n", filename)
	return nil
}

// LoadFighterFromFile loads a fighter object from a JSON file
func LoadFighterFromFile(filename string) (*Fighter, error) {
	// Read the JSON data from the file
	fighterJSON, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error reading fighter data from file: %s", err)
	}

	// Convert the JSON data to a fighter object
	fighter := &Fighter{}
	err = json.Unmarshal(fighterJSON, fighter)
	if err != nil {
		return nil, fmt.Errorf("error decoding fighter from JSON: %s", err)
	}

	fmt.Printf("Fighter data loaded from %s!\n", filename)
	return fighter, nil
}
