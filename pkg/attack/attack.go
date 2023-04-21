package attack

// AttackType represents the type of a fighting move
type AttackType int

const (
	Punch AttackType = iota
	HandStrike
	Kick
	KneeStrike
	ElbowStrike
	Throw
	Lock
	Choke
)

// String returns the string representation of the attack type
func (at AttackType) String() string {
	switch at {
	case Punch:
		return "Punch"
	case HandStrike:
		return "Hand strike"
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
	default:
		return ""
	}
}

// Attack represents an attack in the game
type Attack struct {
	Name           string
	Type           AttackType
	Damage         int
	Complexity     int
	HitChance      int
	BlockChance    int
	CriticalChance int
	SpecialChance  int
}
