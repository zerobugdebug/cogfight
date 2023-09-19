package attack

import (
	"math/rand"

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
	defaultAttacks.AddAttack(&Attack{"Triangle Choke", Choke, 23, 10, 43.5, 12, 50, 90})
	defaultAttacks.AddAttack(&Attack{"Anaconda Choke", Choke, 20, 9.5, 42, 13.5, 50, 90})
	defaultAttacks.AddAttack(&Attack{"D'Arce Choke", Choke, 21, 9, 45, 12.5, 50, 90})
	defaultAttacks.AddAttack(&Attack{"Guillotine Choke", Choke, 22, 8.5, 46.5, 13, 50, 90})
	defaultAttacks.AddAttack(&Attack{"Rear Naked Choke", Choke, 24, 8, 48, 11.5, 50, 90})
	defaultAttacks.AddAttack(&Attack{"Elbow Strike", ElbowStrike, 70, 9, 105, 32.5, 60, 80})
	defaultAttacks.AddAttack(&Attack{"Spinning Heel Kick", Kick, 85, 7.5, 84, 26.5, 72, 15})
	defaultAttacks.AddAttack(&Attack{"Spinning Back Kick", Kick, 75, 7, 87, 27.5, 70, 15})
	defaultAttacks.AddAttack(&Attack{"Roundhouse Kick", Kick, 75, 7, 117, 37.5, 60, 15})
	defaultAttacks.AddAttack(&Attack{"Side Kick", Kick, 65, 6.5, 105, 36.5, 54, 15})
	defaultAttacks.AddAttack(&Attack{"Front Kick", Kick, 65, 6.5, 123, 42.5, 48, 15})
	defaultAttacks.AddAttack(&Attack{"Hook Kick", Kick, 55, 6, 90, 32.5, 50, 15})
	defaultAttacks.AddAttack(&Attack{"Crescent Kick", Kick, 55, 5.5, 96, 38.5, 40, 15})
	defaultAttacks.AddAttack(&Attack{"Axe Kick", Kick, 55, 6, 93, 33.5, 52, 15})
	defaultAttacks.AddAttack(&Attack{"Push Kick", Kick, 15, 2.5, 120, 43.5, 46, 15})
	defaultAttacks.AddAttack(&Attack{"Flying Knee", KneeStrike, 85, 12.5, 90, 27.5, 70, 80})
	defaultAttacks.AddAttack(&Attack{"Knee Strike", KneeStrike, 75, 10, 120, 32.5, 60, 80})
	defaultAttacks.AddAttack(&Attack{"Armbar", Lock, 10, 8.5, 52.5, 19.5, 34, 90})
	defaultAttacks.AddAttack(&Attack{"Kimura", Lock, 10, 8.5, 51, 19, 32, 90})
	defaultAttacks.AddAttack(&Attack{"Americana", Lock, 10, 8.5, 49.5, 18.5, 30, 90})
	defaultAttacks.AddAttack(&Attack{"Straight Foot Lock", Lock, 10, 8.5, 48, 17, 26, 90})
	defaultAttacks.AddAttack(&Attack{"Leg Lock", Lock, 10, 8.5, 46.5, 17.5, 28, 90})
	defaultAttacks.AddAttack(&Attack{"Heel Hook", Lock, 10, 8, 45, 16.5, 62, 90})
	defaultAttacks.AddAttack(&Attack{"Toe Hold", Lock, 10, 8, 43.5, 16, 60, 90})
	defaultAttacks.AddAttack(&Attack{"Omoplata", Lock, 10, 8, 42, 18, 24, 90})
	defaultAttacks.AddAttack(&Attack{"Gogoplata", Lock, 10, 10, 30, 8.5, 22, 90})
	defaultAttacks.AddAttack(&Attack{"Flying Armbar", Lock, 10, 10, 30, 7.5, 20, 90})
	defaultAttacks.AddAttack(&Attack{"Superman Punch", Punch, 45, 2.5, 112.5, 27.5, 62, 10})
	defaultAttacks.AddAttack(&Attack{"Uppercut", Punch, 35, 2, 117, 40.5, 70, 10})
	defaultAttacks.AddAttack(&Attack{"Hook", Punch, 30, 2, 120, 41.5, 58, 10})
	defaultAttacks.AddAttack(&Attack{"Cross", Punch, 25, 1.5, 124.5, 42.5, 52, 10})
	defaultAttacks.AddAttack(&Attack{"Jab", Punch, 15, 1, 127.5, 43.5, 48, 10})
	defaultAttacks.AddAttack(&Attack{"Spinning Back Fist", Slap, 15, 1, 90, 32.5, 70, 100})
	defaultAttacks.AddAttack(&Attack{"Back Fist", Slap, 15, 1, 105, 35.5, 52, 100})
	defaultAttacks.AddAttack(&Attack{"Palm Heel Strike", Slap, 15, 1, 111, 39.5, 48, 100})
	defaultAttacks.AddAttack(&Attack{"Hammer Fist", Slap, 20, 2, 108, 38.5, 50, 100})
	defaultAttacks.AddAttack(&Attack{"Uchi Mata", Throw, 45, 7.5, 66, 17.5, 38, 90})
	defaultAttacks.AddAttack(&Attack{"Harai Goshi", Throw, 45, 7.5, 63, 18.5, 36, 90})
	defaultAttacks.AddAttack(&Attack{"Tai Otoshi", Throw, 45, 7.5, 69, 20.5, 32, 90})
	defaultAttacks.AddAttack(&Attack{"Seoi Nage", Throw, 42, 7, 78, 22.5, 34, 90})
	defaultAttacks.AddAttack(&Attack{"Osoto Gari", Throw, 42, 7, 75, 21.5, 30, 90})
	defaultAttacks.AddAttack(&Attack{"Suplex", Throw, 42, 7.5, 60, 16.5, 40, 90})
	defaultAttacks.AddAttack(&Attack{"Fireman's Carry", Throw, 42, 7, 72, 19.5, 28, 90})
	defaultAttacks.AddAttack(&Attack{"Hip Throw", Throw, 42, 6.5, 87, 24.5, 26, 90})
	defaultAttacks.AddAttack(&Attack{"Shoulder Throw", Throw, 40, 6.5, 84, 23.5, 24, 90})
	defaultAttacks.AddAttack(&Attack{"Foot Sweep", Throw, 40, 6.5, 81, 25.5, 18, 90})
	defaultAttacks.AddAttack(&Attack{"Double Leg Takedown", Throw, 40, 5, 120, 27.5, 22, 90})
	defaultAttacks.AddAttack(&Attack{"Single Leg Takedown", Throw, 40, 5, 117, 26.5, 20, 90})
	defaultAttacks.AddAttack(&Attack{"Ridge Hand", VitalStrike, 15, 13.5, 102, 36.5, 54, 90})
	defaultAttacks.AddAttack(&Attack{"Spear Hand", VitalStrike, 15, 12.5, 99, 37.5, 46, 90})

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
