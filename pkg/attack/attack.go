package attack

import (
	"encoding/csv"
	"math/rand"
	"os"
	"strconv"

	"github.com/zerobugdebug/cogfight/pkg/log"
	"github.com/zerobugdebug/cogfight/pkg/modifiers"
)

const (
	MaxHitChance         float64 = 99
	MinHitChance                 = 1
	MaxBlockChance               = 95
	MinBlockChance               = 0
	MaxComplexity                = 95
	MinComplexity                = 0
	MaxCriticalHitChance         = 95
	MinCriticalHitChance         = 5
	MaxSpecialChance             = 95
	MinSpecialChance             = 5
	MinDamage                    = 5
	MaxDamage                    = 300
)

const (
	defaultAttacksFile = "default_attacks.csv"
)

// AttackType represents the type of a fighting move
type AttackType int

const (
	Punch AttackType = iota
	Slap
	Kick
	KneeStrike
	ElbowStrike
	Throw
	Lock
	Choke
	VitalStrike
	Custom
)

const (
	MaxAttackTypes int = 10
)

var attackTypeNames = map[AttackType]string{
	Punch:       "Punch",
	Slap:        "Slap",
	Kick:        "Kick",
	KneeStrike:  "Knee Strike",
	ElbowStrike: "Elbow Strike",
	Throw:       "Throw",
	Lock:        "Lock",
	Choke:       "Choke",
	VitalStrike: "Vital Strike",
	Custom:      "Custom",
}

// String returns the string representation of the attack type
func (at AttackType) String() string {
	if name, ok := attackTypeNames[at]; ok {
		return name
	}
	return ""
}

var attackTypeHints = map[AttackType]string{
	Punch:       "Closed fist attacks, high damage, low complexity, high hit chance, high block chance",
	Slap:        "Open fist or back hand attacks, very low damage, low complexity, high hit chance, high block chance",
	Kick:        "Leg attacks, high damage, average complexity, high hit chance, high block chance",
	KneeStrike:  "Attacks with a knee, very high damage, average complexity, high hit chance, average block chance",
	ElbowStrike: "Attacks with an elbow, very high damage, low complexity, high hit chance, high block chance",
	Throw:       "Attacks to knockdown opponent, average damage, average complexity, average hit chance, average block chance, can knockdown opponent",
	Lock:        "Grapple attacks to block joint movement, very low damage, high complexity, low hit chance, low block chance, decrease opponent's hit and block chances",
	Choke:       "Grapple attacks to block airways, low damage, high complexity, low hit chance, low block chance, decrease opponent's damage and increase complexity",
	Custom:      "Custom",
}

// Hint returns the string representation of the attack hint
func (at AttackType) Hint() string {
	if hint, ok := attackTypeHints[at]; ok {
		return hint
	}
	return ""
}

var attackTypeSpecials = map[AttackType]modifiers.Condition{
	Punch:       modifiers.CriticalHit,
	Slap:        modifiers.Insulted,
	Kick:        modifiers.CriticalHit,
	KneeStrike:  modifiers.Bleeding,
	ElbowStrike: modifiers.Bleeding,
	Throw:       modifiers.Prone,
	Lock:        modifiers.Bruised,
	Choke:       modifiers.Disoriented,
	VitalStrike: modifiers.Paralysed,
}

// String returns the modifiers value for the Attack type Special
func (at AttackType) Special() modifiers.Condition {
	if special, ok := attackTypeSpecials[at]; ok {
		return special
	}
	return modifiers.Healthy
}

// Attack represents an attack in the game
type Attack struct {
	Name           string
	Type           AttackType
	Damage         float64
	Complexity     float64
	HitChance      float64
	BlockChance    float64
	CriticalChance float64
	SpecialChance  float64
}

// Attacks represents a structure to hold the attacks
type Attacks struct {
	ByName map[string]*Attack
	ByType [MaxAttackTypes][]*Attack
}

func NewAttacks() *Attacks {
	return &Attacks{
		ByName: make(map[string]*Attack),
	}
}

func NewDefaultAttacks() *Attacks {
	log.Infof("Reading configuration file %s", defaultAttacksFile)
	file, err := os.Open(defaultAttacksFile)
	if err != nil {
		log.Fatalf("Failed to open default attack configuration file: %v", err)
		return nil
	}
	defer file.Close()

	// Read the CSV file
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("Failed to read default attack configuration file: %v", err)
		return nil
	}

	// Create a new Attacks struct
	defaultAttacks := NewAttacks()

	// Generate reverse map to attackTypeNames to match strings from the attack CSV file
	var stringToAttackTypeMap = make(map[string]AttackType)

	for k, v := range attackTypeNames {
		stringToAttackTypeMap[v] = k
	}

	// Skip the header row
	for _, record := range records[1:] {
		attackType, ok := stringToAttackTypeMap[record[1]]
		if !ok {
			log.Infof("Unknown attack type: %v. Defaulting to %s", record[1], attackTypeNames[0])
			continue
		}

		damage, _ := strconv.ParseFloat(record[2], 64)
		complexity, _ := strconv.ParseFloat(record[3], 64)
		hitChance, _ := strconv.ParseFloat(record[4], 64)
		blockChance, _ := strconv.ParseFloat(record[5], 64)
		criticalChance, _ := strconv.ParseFloat(record[6], 64)
		specialChance, _ := strconv.ParseFloat(record[7], 64)

		attack := &Attack{
			Name:           record[0],
			Type:           attackType,
			Damage:         damage,
			Complexity:     complexity,
			HitChance:      hitChance,
			BlockChance:    blockChance,
			CriticalChance: criticalChance,
			SpecialChance:  specialChance,
		}

		defaultAttacks.AddAttack(attack)
	}

	return defaultAttacks

}

/* func NewDefaultAttacks() *Attacks {

	defaultAttacks := NewAttacks()

	defaultAttacks.AddAttack(&Attack{"Spear Hand", VitalStrike, 10, 25, 300, 10, 1, 50})
	defaultAttacks.AddAttack(&Attack{"Ridge Hand", VitalStrike, 10, 25, 300, 10, 1, 50})
	defaultAttacks.AddAttack(&Attack{"Flying Knee", KneeStrike, 240, 15, 120, 20, 1, 55})
	defaultAttacks.AddAttack(&Attack{"Knee Strike", KneeStrike, 240, 15, 120, 20, 1, 55})
	defaultAttacks.AddAttack(&Attack{"Rear Naked Choke", Choke, 60, 15, 140, 5, 1, 70})
	defaultAttacks.AddAttack(&Attack{"Guillotine Choke", Choke, 60, 15, 140, 5, 1, 70})
	defaultAttacks.AddAttack(&Attack{"Triangle Choke", Choke, 60, 15, 140, 5, 1, 70})
	defaultAttacks.AddAttack(&Attack{"Anaconda Choke", Choke, 60, 15, 140, 5, 1, 70})
	defaultAttacks.AddAttack(&Attack{"D'Arce Choke", Choke, 60, 15, 140, 5, 1, 70})
	defaultAttacks.AddAttack(&Attack{"Elbow Strike", ElbowStrike, 200, 10, 120, 20, 1, 40})
	defaultAttacks.AddAttack(&Attack{"Front Kick", Kick, 120, 10, 200, 15, 1, 30})
	defaultAttacks.AddAttack(&Attack{"Roundhouse Kick", Kick, 120, 10, 200, 15, 1, 30})
	defaultAttacks.AddAttack(&Attack{"Side Kick", Kick, 120, 10, 200, 15, 1, 30})
	defaultAttacks.AddAttack(&Attack{"Spinning Back Kick", Kick, 120, 10, 200, 15, 1, 30})
	defaultAttacks.AddAttack(&Attack{"Axe Kick", Kick, 120, 10, 200, 15, 1, 30})
	defaultAttacks.AddAttack(&Attack{"Hook Kick", Kick, 120, 10, 200, 15, 1, 30})
	defaultAttacks.AddAttack(&Attack{"Crescent Kick", Kick, 120, 10, 200, 15, 1, 30})
	defaultAttacks.AddAttack(&Attack{"Spinning Heel Kick", Kick, 120, 10, 200, 15, 1, 30})
	defaultAttacks.AddAttack(&Attack{"Push Kick", Kick, 120, 10, 200, 15, 1, 30})
	defaultAttacks.AddAttack(&Attack{"Hip Throw", Throw, 80, 10, 100, 10, 1, 75})
	defaultAttacks.AddAttack(&Attack{"Shoulder Throw", Throw, 80, 10, 100, 10, 1, 75})
	defaultAttacks.AddAttack(&Attack{"Foot Sweep", Throw, 80, 10, 100, 10, 1, 75})
	defaultAttacks.AddAttack(&Attack{"Osoto Gari", Throw, 80, 10, 100, 10, 1, 75})
	defaultAttacks.AddAttack(&Attack{"Uchi Mata", Throw, 80, 10, 100, 10, 1, 75})
	defaultAttacks.AddAttack(&Attack{"Seoi Nage", Throw, 80, 10, 100, 10, 1, 75})
	defaultAttacks.AddAttack(&Attack{"Tai Otoshi", Throw, 80, 10, 100, 10, 1, 75})
	defaultAttacks.AddAttack(&Attack{"Harai Goshi", Throw, 80, 10, 100, 10, 1, 75})
	defaultAttacks.AddAttack(&Attack{"Double Leg Takedown", Throw, 80, 10, 100, 10, 1, 75})
	defaultAttacks.AddAttack(&Attack{"Single Leg Takedown", Throw, 80, 10, 100, 10, 1, 75})
	defaultAttacks.AddAttack(&Attack{"Fireman's Carry", Throw, 80, 10, 100, 10, 1, 75})
	defaultAttacks.AddAttack(&Attack{"Suplex", Throw, 80, 10, 100, 10, 1, 75})
	defaultAttacks.AddAttack(&Attack{"Armbar", Lock, 20, 10, 140, 10, 1, 75})
	defaultAttacks.AddAttack(&Attack{"Kimura", Lock, 20, 10, 140, 10, 1, 75})
	defaultAttacks.AddAttack(&Attack{"Americana", Lock, 20, 10, 140, 10, 1, 75})
	defaultAttacks.AddAttack(&Attack{"Omoplata", Lock, 20, 10, 140, 10, 1, 75})
	defaultAttacks.AddAttack(&Attack{"Gogoplata", Lock, 20, 10, 140, 10, 1, 75})
	defaultAttacks.AddAttack(&Attack{"Leg Lock", Lock, 20, 10, 140, 10, 1, 75})
	defaultAttacks.AddAttack(&Attack{"Heel Hook", Lock, 20, 10, 140, 10, 1, 75})
	defaultAttacks.AddAttack(&Attack{"Straight Foot Lock", Lock, 20, 10, 140, 10, 1, 75})
	defaultAttacks.AddAttack(&Attack{"Toe Hold", Lock, 20, 10, 140, 10, 1, 75})
	defaultAttacks.AddAttack(&Attack{"Flying Armbar", Lock, 20, 10, 140, 10, 1, 75})
	defaultAttacks.AddAttack(&Attack{"Jab", Punch, 100, 5, 280, 15, 1, 10})
	defaultAttacks.AddAttack(&Attack{"Cross", Punch, 100, 5, 280, 15, 1, 10})
	defaultAttacks.AddAttack(&Attack{"Hook", Punch, 100, 5, 280, 15, 1, 10})
	defaultAttacks.AddAttack(&Attack{"Uppercut", Punch, 100, 5, 280, 15, 1, 10})
	defaultAttacks.AddAttack(&Attack{"Superman Punch", Punch, 100, 5, 280, 15, 1, 10})
	defaultAttacks.AddAttack(&Attack{"Palm Heel Strike", Slap, 30, 5, 360, 25, 1, 90})
	defaultAttacks.AddAttack(&Attack{"Hammer Fist", Slap, 30, 5, 360, 25, 1, 90})
	defaultAttacks.AddAttack(&Attack{"Back Fist", Slap, 30, 5, 360, 25, 1, 90})
	defaultAttacks.AddAttack(&Attack{"Spinning Back Fist", Slap, 30, 5, 360, 25, 1, 90})

	return defaultAttacks

} */

func (attacks *Attacks) AddAttack(attack *Attack) {
	// Add attack to the map for lookup by name
	attacks.ByName[attack.Name] = attack

	// Add attack to the corresponding slice for lookup by type
	attacks.ByType[attack.Type] = append(attacks.ByType[attack.Type], attack)
}

func (attacks *Attacks) GetAttackByName(name string) *Attack {
	return attacks.ByName[name]
}

func (attacks *Attacks) GetAttacksByType(attackType AttackType) []*Attack {
	return attacks.ByType[attackType]
}

func (attacks *Attacks) GetRandomAttack() *Attack {
	attackType := AttackType(rand.Intn(MaxAttackTypes - 1))
	attacksNum := len(attacks.GetAttacksByType(attackType))
	return attacks.GetAttacksByType(attackType)[rand.Intn(attacksNum)]
}

func Clamp(val, min, max float64) float64 {
	if val < min {
		return min
	} else if val > max {
		return max
	} else {
		return val
	}
}

/*
type Settable interface {
	WriteAnswer(field string, value interface{}) error
}
*/
