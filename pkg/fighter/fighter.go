package fighter

import (
	"fmt"
	"math/rand"
	"time"
)

// Fighter represents a fighter in the game
type Fighter struct {
	Name         string
	Height       int
	Weight       int
	Age          int
	Speed        int
	Attacks      []string
	CurrentHealth int
	MaxHealth    int
}

// CreateFighter creates a new fighter object based on user input
func CreateFighter() *Fighter {
	// Collect user input
	var name string
	fmt.Print("Enter fighter name: ")
	fmt.Scanln(&name)

	var height int
	fmt.Print("Enter fighter height: ")
	fmt.Scanln(&height)

	var weight int
	fmt.Print("Enter fighter weight: ")
	fmt.Scanln(&weight)

	var age int
	fmt.Print("Enter fighter age: ")
	fmt.Scanln(&age)

	var speed int
	fmt.Print("Enter fighter speed: ")
	fmt.Scanln(&speed)

	attacks := []string{"punch", "kick", "headbutt"}

	// Create the fighter object
	fighter := &Fighter{
		Name:         name,
		Height:       height,
		Weight:       weight,
		Age:          age,
		Speed:        speed,
		Attacks:      attacks,
		CurrentHealth: 100,
		MaxHealth:    100,
	}

	fmt.Printf("\n%s has been created!\n", fighter.Name)
	return fighter
}

// GenerateComputerFighter generates a computer-controlled fighter
func GenerateComputerFighter(playerFighter *Fighter) *Fighter {
	rand.Seed(time.Now().UnixNano())

	// Generate random stats for the computer fighter
	name := "Computer Fighter"
	height := rand.Intn(30) + 160 // Computer fighter's height is between 160cm to 190cm
	weight := rand.Intn(30) + 60 // Computer fighter's weight is between 60kg to 90kg
	age := rand.Intn(20) + 20 // Computer fighter's age is between 20 to 40 years
	speed := rand.Intn(10) + 20 // Computer fighter's speed is between 20 to 30 m/s

	attacks := []string{"punch", "kick", "headbutt"}

	// Create the computer fighter object
	computerFighter := &Fighter{
		Name:         name,
		Height:       height,
		Weight:       weight,
		Age:          age,
		Speed:        speed,
		Attacks:      attacks,
		CurrentHealth: 100,
		MaxHealth:    100,
	}

	fmt.Printf("\n%s has been created as your opponent!\n", computerFighter.Name)
	return computerFighter
}