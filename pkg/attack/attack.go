package attack

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

// String returns the string representation of the attack type
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

// Attack represents an attack in the game
type Attack struct {
	Name           string
	Type           AttackType
	Damage         float32
	Complexity     float32
	HitChance      float32
	BlockChance    float32
	CriticalChance float32
	SpecialChance  float32
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
	defaultAttacks.AddAttack(&Attack{"Jab", Punch, 15, 14, 80, 87, 48, 15})
	defaultAttacks.AddAttack(&Attack{"Cross", Punch, 55, 15, 78, 85, 52, 15})
	defaultAttacks.AddAttack(&Attack{"Hook", Punch, 65, 16, 70, 83, 58, 15})
	defaultAttacks.AddAttack(&Attack{"Uppercut", Punch, 65, 17, 68, 81, 70, 15})
	defaultAttacks.AddAttack(&Attack{"Front Kick", Kick, 65, 25, 82, 85, 48, 45})
	defaultAttacks.AddAttack(&Attack{"Roundhouse Kick", Kick, 75, 27, 78, 75, 60, 45})
	defaultAttacks.AddAttack(&Attack{"Side Kick", Kick, 65, 29, 70, 73, 54, 30})
	defaultAttacks.AddAttack(&Attack{"Spinning Back Kick", Kick, 75, 45, 58, 55, 70, 45})
	defaultAttacks.AddAttack(&Attack{"Axe Kick", Kick, 55, 33, 62, 67, 52, 30})
	defaultAttacks.AddAttack(&Attack{"Hook Kick", Kick, 55, 37, 60, 65, 50, 15})
	defaultAttacks.AddAttack(&Attack{"Crescent Kick", Kick, 55, 35, 64, 77, 40, 15})
	defaultAttacks.AddAttack(&Attack{"Spinning Heel Kick", Kick, 85, 47, 56, 53, 72, 15})
	defaultAttacks.AddAttack(&Attack{"Superman Punch", Punch, 85, 45, 60, 55, 62, 30})
	defaultAttacks.AddAttack(&Attack{"Flying Knee", KneeStrike, 85, 85, 60, 55, 70, 45})
	defaultAttacks.AddAttack(&Attack{"Elbow Strike", ElbowStrike, 75, 35, 70, 65, 60, 15})
	defaultAttacks.AddAttack(&Attack{"Knee Strike", KneeStrike, 75, 35, 80, 65, 60, 30})
	defaultAttacks.AddAttack(&Attack{"Palm Heel Strike", Slap, 15, 17, 74, 79, 48, 30})
	defaultAttacks.AddAttack(&Attack{"Hammer Fist", Slap, 55, 15, 72, 77, 50, 15})
	defaultAttacks.AddAttack(&Attack{"Spear Hand", VitalStrike, 15, 35, 66, 75, 46, 15})
	defaultAttacks.AddAttack(&Attack{"Ridge Hand", VitalStrike, 15, 37, 68, 73, 54, 15})
	defaultAttacks.AddAttack(&Attack{"Back Fist", Slap, 15, 19, 70, 71, 52, 15})
	defaultAttacks.AddAttack(&Attack{"Push Kick", Kick, 15, 23, 80, 87, 46, 45})
	defaultAttacks.AddAttack(&Attack{"Spinning Back Fist", Slap, 15, 45, 60, 65, 70, 15})
	defaultAttacks.AddAttack(&Attack{"Hip Throw", Throw, 45, 51, 58, 49, 26, 75})
	defaultAttacks.AddAttack(&Attack{"Shoulder Throw", Throw, 45, 49, 56, 47, 24, 75})
	defaultAttacks.AddAttack(&Attack{"Foot Sweep", Throw, 45, 47, 54, 51, 18, 75})
	defaultAttacks.AddAttack(&Attack{"Osoto Gari", Throw, 45, 57, 50, 43, 30, 75})
	defaultAttacks.AddAttack(&Attack{"Uchi Mata", Throw, 45, 65, 44, 35, 38, 75})
	defaultAttacks.AddAttack(&Attack{"Seoi Nage", Throw, 45, 59, 52, 45, 34, 75})
	defaultAttacks.AddAttack(&Attack{"Tai Otoshi", Throw, 45, 61, 46, 41, 32, 75})
	defaultAttacks.AddAttack(&Attack{"Harai Goshi", Throw, 45, 63, 42, 37, 36, 75})
	defaultAttacks.AddAttack(&Attack{"Double Leg Takedown", Throw, 45, 45, 80, 55, 22, 75})
	defaultAttacks.AddAttack(&Attack{"Single Leg Takedown", Throw, 45, 43, 78, 53, 20, 75})
	defaultAttacks.AddAttack(&Attack{"Fireman's Carry", Throw, 45, 53, 48, 39, 28, 60})
	defaultAttacks.AddAttack(&Attack{"Suplex", Throw, 45, 55, 40, 33, 40, 60})
	defaultAttacks.AddAttack(&Attack{"Rear Naked Choke", Choke, 37, 64, 32, 23, 50, 50})
	defaultAttacks.AddAttack(&Attack{"Guillotine Choke", Choke, 35, 65, 31, 26, 50, 50})
	defaultAttacks.AddAttack(&Attack{"Triangle Choke", Choke, 36, 75, 29, 24, 50, 50})
	defaultAttacks.AddAttack(&Attack{"Armbar", Lock, 25, 63, 35, 39, 34, 50})
	defaultAttacks.AddAttack(&Attack{"Kimura", Lock, 25, 69, 34, 38, 32, 50})
	defaultAttacks.AddAttack(&Attack{"Americana", Lock, 25, 67, 33, 37, 30, 50})
	defaultAttacks.AddAttack(&Attack{"Omoplata", Lock, 25, 77, 28, 36, 24, 50})
	defaultAttacks.AddAttack(&Attack{"Gogoplata", Lock, 25, 85, 20, 17, 22, 50})
	defaultAttacks.AddAttack(&Attack{"Leg Lock", Lock, 25, 71, 31, 35, 28, 50})
	defaultAttacks.AddAttack(&Attack{"Anaconda Choke", Choke, 33, 67, 28, 27, 50, 50})
	defaultAttacks.AddAttack(&Attack{"D'Arce Choke", Choke, 34, 66, 30, 25, 50, 50})
	defaultAttacks.AddAttack(&Attack{"Heel Hook", Lock, 25, 75, 30, 33, 62, 50})
	defaultAttacks.AddAttack(&Attack{"Straight Foot Lock", Lock, 25, 65, 32, 34, 26, 50})
	defaultAttacks.AddAttack(&Attack{"Toe Hold", Lock, 25, 73, 29, 32, 60, 50})
	defaultAttacks.AddAttack(&Attack{"Flying Armbar", Lock, 25, 83, 20, 15, 20, 50})
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

/*
type Settable interface {
	WriteAnswer(field string, value interface{}) error
}
*/
