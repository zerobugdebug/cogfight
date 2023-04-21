package fighter

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/sashabaranov/go-openai"
	"github.com/zerobugdebug/cogfight/pkg/attack"
)

const (
	numAttacks int = 3
)

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
	SpecialChanceBonus          float32
	Attacks                     []*attack.Attack
	CurrentHealth               int
	MaxHealth                   int
}

func (f *Fighter) DisplayFighter() {
	topBorder := "╔══════════════════════════════════════════════════════════╗"
	bottomBorder := "╚══════════════════════════════════════════════════════════╝"
	spacer := "║                                                          ║"

	fmt.Println(topBorder)
	fmt.Println(spacer)
	fmt.Printf("║ Name: %-50s ║\n", f.Name)
	fmt.Printf("║ Height: %-48d ║\n", f.Height)
	fmt.Printf("║ Weight: %-48d ║\n", f.Weight)
	fmt.Printf("║ Age: %-51d ║\n", f.Age)
	fmt.Printf("║ Agility Strength Balance: %-30.2f ║\n", f.AgilityStrengthBalance)
	fmt.Printf("║ Burst Endurance Balance: %-31.2f ║\n", f.BurstEnduranceBalance)
	fmt.Printf("║ Defense Offense Balance: %-31.2f ║\n", f.DefenseOffenseBalance)
	fmt.Printf("║ Speed Control Balance: %-33.2f ║\n", f.SpeedControlBalance)
	fmt.Printf("║ Intelligence Instinct Balance: %-25.2f ║\n", f.IntelligenceInstinctBalance)
	fmt.Printf("║ Damage Bonus: %-42.2f ║\n", f.DamageBonus)
	fmt.Printf("║ Complexity Bonus: %-38.2f ║\n", f.ComplexityBonus)
	fmt.Printf("║ Hit Chance Bonus: %-38.2f ║\n", f.HitChanceBonus)
	fmt.Printf("║ Block Chance Bonus: %-36.2f ║\n", f.BlockChanceBonus)
	fmt.Printf("║ Critical Chance Bonus: %-33.2f ║\n", f.CriticalChanceBonus)
	fmt.Printf("║ Special Chance Bonus: %-34.2f ║\n", f.SpecialChanceBonus)
	fmt.Println(spacer)
	fmt.Println("║ Attacks:                                                  ║")
	for i, attack := range f.Attacks {
		fmt.Printf("║ %d. %-53s ║\n", i+1, attack.Name)
	}
	fmt.Println(spacer)
	fmt.Printf("║ Current Health: %-38d ║\n", f.CurrentHealth)
	fmt.Printf("║ Max Health: %-41d ║\n", f.MaxHealth)
	fmt.Println(spacer)
	fmt.Println(bottomBorder)
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
			Help:    "Please enter your fighter height. Taller fighters will favour Strength, Offense and Control, while lower height will give benefits to Agility, Defense and Speed.",
		},
		Validate: validateNumber(150, 220),
	}
	qs = append(qs, heightQuestion)

	weightQuestion := &survey.Question{
		Name: "weight",
		Prompt: &survey.Input{
			Message: "Enter fighter weight (50-200 kg):",
			Help:    "Please enter your fighter weight. Heavier fighters tend to have better Strength, Endurance and Control, while lighter fighters rely more on the Agility, Burst and Speed.",
		},
		Validate: validateNumber(50, 200),
	}
	qs = append(qs, weightQuestion)

	ageQuestion := &survey.Question{
		Name: "age",
		Prompt: &survey.Input{
			Message: "Enter fighter age (18-60 years):",
			Help:    "Please enter your fighter age. Older fighters tend to have better Intelligence, while younger fighters rely more on the Instinct.",
		},
		Validate: validateNumber(18, 60),
	}
	qs = append(qs, ageQuestion)

	agilityStrengthBalanceQuestion := &survey.Question{
		Name: "agilityStrengthBalance",
		Prompt: &survey.Select{
			Message: "Choose fighter agility/strength balance:",
			Options: []string{"Very high Agility, Very low Strength", "High Agility, Low Strength", "Balanced", "Low Agility, High Strength", "Very low Agility, Very high Strength"},
			Help:    "This parameter determines the balance between Agility and Strength. High Agility will allow fighter to execute more complex attack with better chances to hit and block, while high Strength will increase damage and special effects chance.",
		},
	}
	qs = append(qs, agilityStrengthBalanceQuestion)

	burstEnduranceBalanceQuestion := &survey.Question{
		Name: "burstEnduranceBalance",
		Prompt: &survey.Select{
			Message: "Choose fighter burst/endurance balance:",
			Options: []string{"Very high Burst, Very low Endurance", "High Burst, Low Endurance", "Balanced", "Low Burst, High Endurance", "Very low Burst, Very high Endurance"},
			Help:    "This parameter determines the balance between Burst and Endurance. Fighters with high Burst will get better chances to hit and special effects, but high Endurance will give bonuses to damage and blocking chance.",
		},
	}
	qs = append(qs, burstEnduranceBalanceQuestion)

	defenseOffenseBalanceQuestion := &survey.Question{
		Name: "defenseOffenseBalance",
		Prompt: &survey.Select{
			Message: "Choose fighter defense/offense balance:",
			Options: []string{"Very high Defense, Very low Offense", "High Defense, Low Offense", "Balanced", "Low Defense, High Offense", "Very low Defense, Very high Offense"},
			Help:    "This parameter determines the balance between Defense and Offense. Increasing Defense will improve your chances of blocking attacks, while increasing Offense will help with hitting.",
		},
	}
	qs = append(qs, defenseOffenseBalanceQuestion)

	speedControlBalanceQuestion := &survey.Question{
		Name: "speedControlBalance",
		Prompt: &survey.Select{
			Message: "Choose fighter speed/control balance:",
			Options: []string{"Very high Speed, Very low Control", "High Speed, Low Control", "Balanced", "Low Speed, High Control", "Very low Speed, Very high Control"},
			Help:    "This parameter determines the balance between Speed and Control. Increasing Speed will improve your chances of successfully hitting and blocking attacks, while high Control will help with executing more complex attacks and critical hits",
		},
	}
	qs = append(qs, speedControlBalanceQuestion)

	intelligenceInstinctBalanceQuestion := &survey.Question{
		Name: "intelligenceInstinctBalance",
		Prompt: &survey.Select{
			Message: "Choose fighter intelligence/instinct balance:",
			Options: []string{"Very high Intelligence, Very low Instinct", "High Intelligence, Low Instinct", "Balanced", "Low Intelligence, High Instinct", "Very low Intelligence, Very high Instinct"},
			Help:    "This parameter determines the balance between Intelligence and Instinct. Increasing Intelligence will help with executing more complex attacks and critical hits, while increasing Instinct will improve your chances of successfully hitting and blocking attacks",
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

	attackType := attack.AttackType(0)
	attackTypePromptOptions := []string{}

	for attackType.String() != "" {
		attackTypePromptOptions = append(attackTypePromptOptions, attackType.String())
		attackType++
	}

	attackTypePrompt := &survey.Select{
		Message:  "Select an attack type:",
		Options:  attackTypePromptOptions,
		PageSize: attack.MaxAttackTypes,
	}

	defaultAttacks := attack.NewDefaultAttacks()
	attacks := []*attack.Attack{}

	i := 0
	for i < numAttacks {
		// TODO: Add number of attack in the question, like 1 of 3
		attackTypeSelected := 0
		// Ask for attack type
		err := survey.AskOne(attackTypePrompt, &attackTypeSelected, survey.WithValidator(survey.Required))
		if err != nil {
			fmt.Println(err)
			break
		}
		attackType = attack.AttackType(attackTypeSelected)
		fmt.Println(defaultAttacks.GetAttacksByType(attackType))

		// If non-custom type, ask for specific attack
		if attackType != attack.Custom {
			attackNamePromptOptions := []string{}
			for _, value := range defaultAttacks.GetAttacksByType(attackType) {
				attackNamePromptOptions = append(attackNamePromptOptions, value.Name)
			}
			attackNamePrompt := &survey.Select{
				Message:  "Select an attack:",
				Options:  attackNamePromptOptions,
				PageSize: len(attackNamePromptOptions),
			}
			attackName := ""
			err = survey.AskOne(attackNamePrompt, &attackName, survey.WithValidator(survey.Required))
			if err != nil {
				fmt.Println(err)
				break
			}
			// Add attack to attacks array
			attacks = append(attacks, defaultAttacks.GetAttackByName(attackName))
			i++
			continue
		}

		// Ask for the description of the custom attack
		customAttackName := ""
		customAttackPrompt := &survey.Input{
			Message: "Enter the description for the custom attack:",
		}
		err = survey.AskOne(customAttackPrompt, &customAttackName)
		if err != nil {
			fmt.Println(err)
			break
		}

		// Validate the attack name and get the attack parameters using OpenAI API
		validAttack, err := validateAttackName(customAttackName)
		if err != nil {
			fmt.Printf("Error validating attack name: %s\n", err)
			continue
		}

		if validAttack {
			complexityValue, err := getIntegerOpenAIResponse("COG_COMPLEXITY_ATTACK_PROMPT", customAttackName)
			if err != nil {
				fmt.Printf("Error getting data for COG_COMPLEXITY_ATTACK_PROMPT: %s\n", err)
				continue
			}
			fmt.Println("complexityValue=", complexityValue)
		}
	}

	// Create the fighter object
	fighter := &Fighter{
		Name:                        answers.Name,
		Height:                      answers.Height,
		Weight:                      answers.Weight,
		Age:                         answers.Age,
		AgilityStrengthBalance:      float32(answers.AgilityStrengthBalance) + (float32(answers.Weight)-125)/50 + (float32(answers.Height)-185)/20,
		BurstEnduranceBalance:       float32(answers.BurstEnduranceBalance) + (float32(answers.Weight)-125)/50,
		DefenseOffenseBalance:       float32(answers.DefenseOffenseBalance) + (float32(answers.Height)-185)/20,
		SpeedControlBalance:         float32(answers.SpeedControlBalance) + (float32(answers.Weight)-125)/50 + (float32(answers.Height)-185)/20,
		IntelligenceInstinctBalance: float32(answers.IntelligenceInstinctBalance) - (float32(answers.Age)-39)/10,
		CurrentHealth:               100,
		MaxHealth:                   100,
		Attacks:                     attacks,
	}

	fighter.DamageBonus = (fighter.AgilityStrengthBalance + fighter.BurstEnduranceBalance - 4) * 10
	fighter.ComplexityBonus = (-fighter.AgilityStrengthBalance + fighter.SpeedControlBalance - fighter.IntelligenceInstinctBalance + 2) * 2
	fighter.HitChanceBonus = (-fighter.AgilityStrengthBalance - fighter.BurstEnduranceBalance + fighter.DefenseOffenseBalance - fighter.SpeedControlBalance + fighter.IntelligenceInstinctBalance + 2) * 2
	fighter.BlockChanceBonus = (-fighter.AgilityStrengthBalance + fighter.BurstEnduranceBalance - fighter.DefenseOffenseBalance - fighter.SpeedControlBalance + fighter.IntelligenceInstinctBalance + 2) * 2
	fighter.CriticalChanceBonus = (fighter.SpeedControlBalance - fighter.IntelligenceInstinctBalance) * 2
	fighter.SpecialChanceBonus = (fighter.AgilityStrengthBalance - fighter.BurstEnduranceBalance) * 2

	fmt.Printf("\n%s has been created!\n", fighter.Name)
	fighter.DisplayFighter()

	return fighter
}

// GenerateComputerFighter generates a computer-controlled fighter
func GenerateComputerFighter(playerFighter *Fighter) *Fighter {

	answers := struct {
		Height                      int
		Weight                      int
		Age                         int
		AgilityStrengthBalance      int
		BurstEnduranceBalance       int
		DefenseOffenseBalance       int
		SpeedControlBalance         int
		IntelligenceInstinctBalance int
	}{}
	answers.AgilityStrengthBalance = rand.Intn(5)
	answers.BurstEnduranceBalance = rand.Intn(5)
	answers.DefenseOffenseBalance = rand.Intn(5)
	answers.SpeedControlBalance = rand.Intn(5)
	answers.IntelligenceInstinctBalance = rand.Intn(5)

	answers.Height = rand.Intn(71) + 150 // Height between 150 and 199 cm
	answers.Weight = rand.Intn(151) + 50 // Weight between 50 and 99 kg
	answers.Age = rand.Intn(43) + 18     // Age between 20 and 49 years

	// Generate random values for the computer fighter's attributes
	computerFighter := &Fighter{
		Name:                        "Computer Fighter",
		Height:                      answers.Height,
		Weight:                      answers.Weight,
		Age:                         answers.Age,
		AgilityStrengthBalance:      float32(answers.AgilityStrengthBalance) + (float32(answers.Weight)-125)/50 + (float32(answers.Height)-185)/20,
		BurstEnduranceBalance:       float32(answers.BurstEnduranceBalance) + (float32(answers.Weight)-125)/50,
		DefenseOffenseBalance:       float32(answers.DefenseOffenseBalance) + (float32(answers.Height)-185)/20,
		SpeedControlBalance:         float32(answers.SpeedControlBalance) + (float32(answers.Weight)-125)/50 + (float32(answers.Height)-185)/20,
		IntelligenceInstinctBalance: float32(answers.IntelligenceInstinctBalance) - (float32(answers.Age)-39)/10,
		Attacks:                     []*attack.Attack{},
		CurrentHealth:               100,
		MaxHealth:                   100,
	}

	computerFighter.DamageBonus = (computerFighter.AgilityStrengthBalance + computerFighter.BurstEnduranceBalance - 4) * 10
	computerFighter.ComplexityBonus = (-computerFighter.AgilityStrengthBalance + computerFighter.SpeedControlBalance - computerFighter.IntelligenceInstinctBalance + 2) * 2
	computerFighter.HitChanceBonus = (-computerFighter.AgilityStrengthBalance - computerFighter.BurstEnduranceBalance + computerFighter.DefenseOffenseBalance - computerFighter.SpeedControlBalance + computerFighter.IntelligenceInstinctBalance + 2) * 2
	computerFighter.BlockChanceBonus = (-computerFighter.AgilityStrengthBalance + computerFighter.BurstEnduranceBalance - computerFighter.DefenseOffenseBalance - computerFighter.SpeedControlBalance + computerFighter.IntelligenceInstinctBalance + 2) * 2
	computerFighter.CriticalChanceBonus = (computerFighter.SpeedControlBalance - computerFighter.IntelligenceInstinctBalance) * 2
	computerFighter.SpecialChanceBonus = (computerFighter.AgilityStrengthBalance - computerFighter.BurstEnduranceBalance) * 2

	defaultAttacks := attack.NewDefaultAttacks()
	for range playerFighter.Attacks {
		attacksList := defaultAttacks.GetAttacksByType(attack.AttackType(rand.Intn(attack.MaxAttackTypes - 1)))
		fmt.Println(attacksList)
		fmt.Println(len(attacksList))
		computerAttack := attacksList[rand.Intn(len(attacksList))]
		fmt.Println(computerAttack)
		computerFighter.Attacks = append(computerFighter.Attacks, computerAttack)
	}

	fmt.Printf("\n%s has been generated!\n", computerFighter.Name)
	computerFighter.DisplayFighter()
	return computerFighter
}

// validateAttackName validates the given attack name using OpenAI API and returns the attack parameters
func validateAttackName(attackName string) (bool, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return false, fmt.Errorf("OpenAI API key not found in environment variable OPENAI_API_KEY")
	}

	client := openai.NewClient(apiKey)

	// Define the prompt template
	promptTemplate := os.Getenv("COG_VALIDATION_ATTACK_PROMPT")

	// Send the prompt to OpenAI API and get the response
	prompt := fmt.Sprintf(promptTemplate, attackName)

	response, err := client.CreateChatCompletion(
		context.Background(),

		openai.ChatCompletionRequest{
			Model:     openai.GPT4,
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
		return false, fmt.Errorf("error sending OpenAI API request: %s", err)
	}

	fmt.Println(response.Choices[0].Message.Content)
	// Parse the response and extract integer answer
	reply := response.Choices[0].Message.Content

	// Parse the response to confirm if attack is valid
	//reply := response.Choices[0].Message.Content
	client = nil
	fmt.Println(reply)

	if strings.Contains(reply, "Invalid") {
		return false, errors.New("Not a valid attack")
	} else if strings.Contains(reply, "Multiple") {
		return false, errors.New("Attack is valid, but not a single attack")
	} else if strings.Contains(reply, "Impossible") || strings.Contains(reply, "Valid") {
		return true, nil
	}

	return false, fmt.Errorf("Unknown response from OpenAI API: %s", reply)
}

// Get integer answer from OpenAI API
func getIntegerOpenAIResponse(promptEnvVariable string, promptData string) (int, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return 0, fmt.Errorf("OpenAI API key not found in environment variable OPENAI_API_KEY")
	}

	client := openai.NewClient(apiKey)

	promptTemplate := os.Getenv(promptEnvVariable)

	if promptTemplate == "" {
		return 0, fmt.Errorf("Prompt not found in the environment variable %s", promptEnvVariable)
	}

	// Send the prompt to OpenAI API and get the response
	prompt := fmt.Sprintf(promptTemplate, promptData)

	response, err := client.CreateChatCompletion(
		context.Background(),

		openai.ChatCompletionRequest{
			Model:     openai.GPT4,
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
		return 0, fmt.Errorf("Error sending OpenAI API request: %s", err)
	}

	fmt.Println(response.Choices[0].Message.Content)
	// Parse the response and extract integer answer
	reply := response.Choices[0].Message.Content
	re := regexp.MustCompile(`\[\[(\d+)\]\]`)
	match := re.FindStringSubmatch(reply)
	fmt.Println(match)
	if len(match) > 1 {
		fmt.Println("Number:", match[1])
		return strconv.Atoi(match[1])
	}

	return 0, fmt.Errorf("Can't parse OpenAI API response: %s", reply)
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
