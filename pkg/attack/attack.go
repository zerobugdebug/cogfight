package attack

import (
	"math/rand"

	"github.com/zerobugdebug/cogfight/pkg/modifiers"


const (
	MaxHitChance         float64 = 99
	MinHitChance                 = 1
	MaxBlockChance               = 95
	MinBlockChance               = 0
	MaxComplexity                = 95
	MinComplexity                = 0
	MaxCriticalHitChance         = 95
	MinCriticalHitChance         = 5
	MaxSpecialChance             = 100
	MinSpecialChance             = 5
	MinDamage                    = 5
	MaxDamage                    = 100
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

// String returns the string representation of the attack type
func (at AttackType) String() string {
	switch at {
	case Punch:
		return "Punch"
	case Slap:
		return "Slap"
	case Kick:
		return "Kick"
	case KneeStrike:
		return "Knee strike"
	case ElbowStrike:
		return "Elbow strike"
	case Throw:
		return "Throw"
	case Lock:
		return "Lock"
	case Choke:
		return "Choke"
	case VitalStrike:
		return "Vital strike"
	case Custom:
		return "Custom"
	default:
		return ""
	}
}

// Hint returns the string representation of the attack type
func (at AttackType) Hint() string {
	switch at {
	case Punch:
		return "Closed fist attacks, high damage, low complexity, high hit chance, high block chance"
	case Slap:
		return "Open fist or back hand attacks, very low damage, low complexity, high hit chance, high block chance"
	case Kick:
		return "Leg attacks, high damage, average complexity, high hit chance, high block chance"
	case KneeStrike:
		return "Attacks with a knee, very high damage, average complexity, high hit chance, average block chance"
	case ElbowStrike:
		return "Attacks with an elbow, very high damage, low complexity, high hit chance, high block chance"
	case Throw:
		return "Attacks to knockdown opponent, average damage, average complexity, average hit chance, average block chance, can knockdown opponent"
	case Lock:
		return "Grapple attacks to block joint movement, very low damage, high complexity, low hit chance, low block chance, decrease opponent's hit and block chances"
	case Choke:
		return "Grapple attacks to block airways, low damage, high complexity, low hit chance, low block chance, decrease opponent's damage and increase complexity"
	case Custom:
		return "Custom"
	default:
		return ""
	}
}

// String returns the string representation of the attack type
func (at AttackType) Special() modifiers.Condition {
	switch at {
	case Punch:
		return modifiers.CriticalHit
	case Slap:
		return modifiers.Insulted
	case Kick:
		return modifiers.CriticalHit
	case KneeStrike:
		return modifiers.Bleeding
	case ElbowStrike:
		return modifiers.Bleeding
	case Throw:
		return modifiers.Prone
	case Lock:
		return modifiers.Bruised
	case Choke:
		return modifiers.Disoriented
	case VitalStrike:
		return modifiers.Paralysed
	default:
		return modifiers.Healthy
	}
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
	defaultAttacks := NewAttacks()
	defaultAttacks.AddAttack(&Attack{"Triangle Choke", Choke, 30, 25, 15, 10, 1, 70})
	defaultAttacks.AddAttack(&Attack{"Anaconda Choke", Choke, 30, 25, 15, 10, 1, 70})
	defaultAttacks.AddAttack(&Attack{"D'Arce Choke", Choke, 30, 25, 15, 10, 1, 70})
	defaultAttacks.AddAttack(&Attack{"Guillotine Choke", Choke, 30, 25, 15, 10, 1, 70})
	defaultAttacks.AddAttack(&Attack{"Rear Naked Choke", Choke, 30, 25, 15, 10, 1, 70})
	defaultAttacks.AddAttack(&Attack{"Elbow Strike", ElbowStrike, 90, 17.5, 5, 70, 1, 50})
	defaultAttacks.AddAttack(&Attack{"Spinning Heel Kick", Kick, 85, 15, 50, 55, 1, 30})
	defaultAttacks.AddAttack(&Attack{"Spinning Back Kick", Kick, 85, 15, 50, 55, 1, 30})
	defaultAttacks.AddAttack(&Attack{"Roundhouse Kick", Kick, 85, 15, 50, 55, 1, 30})
	defaultAttacks.AddAttack(&Attack{"Side Kick", Kick, 85, 15, 50, 55, 1, 30})
	defaultAttacks.AddAttack(&Attack{"Front Kick", Kick, 85, 15, 50, 55, 1, 30})
	defaultAttacks.AddAttack(&Attack{"Hook Kick", Kick, 85, 15, 50, 55, 1, 30})
	defaultAttacks.AddAttack(&Attack{"Crescent Kick", Kick, 85, 15, 50, 55, 1, 30})
	defaultAttacks.AddAttack(&Attack{"Axe Kick", Kick, 85, 15, 50, 55, 1, 30})
	defaultAttacks.AddAttack(&Attack{"Push Kick", Kick, 85, 15, 50, 55, 1, 30})
	defaultAttacks.AddAttack(&Attack{"Flying Knee", KneeStrike, 95, 30, 30, 95, 1, 65})
	defaultAttacks.AddAttack(&Attack{"Knee Strike", KneeStrike, 95, 30, 30, 95, 1, 65})
	defaultAttacks.AddAttack(&Attack{"Armbar", Lock, 10, 12.5, 35, 15, 1, 75})
	defaultAttacks.AddAttack(&Attack{"Kimura", Lock, 10, 12.5, 35, 15, 1, 75})
	defaultAttacks.AddAttack(&Attack{"Americana", Lock, 10, 12.5, 35, 15, 1, 75})
	defaultAttacks.AddAttack(&Attack{"Straight Foot Lock", Lock, 10, 12.5, 35, 15, 1, 75})
	defaultAttacks.AddAttack(&Attack{"Leg Lock", Lock, 10, 12.5, 35, 15, 1, 75})
	defaultAttacks.AddAttack(&Attack{"Heel Hook", Lock, 10, 12.5, 35, 15, 1, 75})
	defaultAttacks.AddAttack(&Attack{"Toe Hold", Lock, 10, 12.5, 35, 15, 1, 75})
	defaultAttacks.AddAttack(&Attack{"Omoplata", Lock, 10, 12.5, 35, 15, 1, 75})
	defaultAttacks.AddAttack(&Attack{"Gogoplata", Lock, 10, 12.5, 35, 15, 1, 75})
	defaultAttacks.AddAttack(&Attack{"Flying Armbar", Lock, 10, 12.5, 35, 15, 1, 75})
	defaultAttacks.AddAttack(&Attack{"Superman Punch", Punch, 70, 7.5, 70, 50, 1, 10})
	defaultAttacks.AddAttack(&Attack{"Uppercut", Punch, 70, 7.5, 70, 50, 1, 10})
	defaultAttacks.AddAttack(&Attack{"Hook", Punch, 70, 7.5, 70, 50, 1, 10})
	defaultAttacks.AddAttack(&Attack{"Cross", Punch, 70, 7.5, 70, 50, 1, 10})
	defaultAttacks.AddAttack(&Attack{"Jab", Punch, 70, 7.5, 70, 50, 1, 10})
	defaultAttacks.AddAttack(&Attack{"Spinning Back Fist", Slap, 15, 5, 90, 90, 1, 90})
	defaultAttacks.AddAttack(&Attack{"Back Fist", Slap, 15, 5, 90, 90, 1, 90})
	defaultAttacks.AddAttack(&Attack{"Palm Heel Strike", Slap, 15, 5, 90, 90, 1, 90})
	defaultAttacks.AddAttack(&Attack{"Hammer Fist", Slap, 15, 5, 90, 90, 1, 90})
	defaultAttacks.AddAttack(&Attack{"Uchi Mata", Throw, 50, 27.5, 10, 35, 1, 60})
	defaultAttacks.AddAttack(&Attack{"Harai Goshi", Throw, 50, 27.5, 10, 35, 1, 60})
	defaultAttacks.AddAttack(&Attack{"Tai Otoshi", Throw, 50, 27.5, 10, 35, 1, 60})
	defaultAttacks.AddAttack(&Attack{"Seoi Nage", Throw, 50, 27.5, 10, 35, 1, 60})
	defaultAttacks.AddAttack(&Attack{"Osoto Gari", Throw, 50, 27.5, 10, 35, 1, 60})
	defaultAttacks.AddAttack(&Attack{"Suplex", Throw, 50, 27.5, 10, 35, 1, 60})
	defaultAttacks.AddAttack(&Attack{"Fireman's Carry", Throw, 50, 27.5, 10, 35, 1, 60})
	defaultAttacks.AddAttack(&Attack{"Hip Throw", Throw, 50, 27.5, 10, 35, 1, 60})
	defaultAttacks.AddAttack(&Attack{"Shoulder Throw", Throw, 50, 27.5, 10, 35, 1, 60})
	defaultAttacks.AddAttack(&Attack{"Foot Sweep", Throw, 50, 27.5, 10, 35, 1, 60})
	defaultAttacks.AddAttack(&Attack{"Double Leg Takedown", Throw, 50, 27.5, 10, 35, 1, 60})
	defaultAttacks.AddAttack(&Attack{"Single Leg Takedown", Throw, 50, 27.5, 10, 35, 1, 60})
	defaultAttacks.AddAttack(&Attack{"Ridge Hand", VitalStrike, 5, 45, 75, 30, 1, 50})
	defaultAttacks.AddAttack(&Attack{"Spear Hand", VitalStrike, 5, 45, 75, 30, 1, 50})
		
	return defaultAttacks

}

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
