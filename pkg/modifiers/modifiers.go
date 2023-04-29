package modifiers

// Condition is an enumeration of fighter conditions
type Condition int

const (
	Healthy Condition = iota
	Bruised
	Disoriented
	Prone
	CriticalHit
	Bleeding
	Paralysed
	Insulted
)

type Modifier int

const (
	HitChance Modifier = iota
	BlockChance
	Duration
	Damage
	DamageMult
	Complexity
	SkipTurn
	OpponentHitChance
	OpponentBlockChance
	HPPerTurn
	MyHPPerTurn
	SureStrike
)

// ConditionAttributes maps fighter conditions to their respective attributes
var DefaultConditionAttributes = map[Condition]map[Modifier]int{
	Healthy: {},
	Bruised: {
		HitChance:   -20,
		BlockChance: -20,
		Duration:    3,
	},
	Disoriented: {
		Damage:     -20,
		Complexity: 20,
		Duration:   3,
	},
	Prone: {
		SkipTurn: 1,
		Duration: 1,
	},
	CriticalHit: {
		DamageMult: 5,
		Duration:   1,
	},
	Bleeding: {
		HPPerTurn: -20,
		Duration:  3,
	},
	Paralysed: {
		SureStrike: 1,
		SkipTurn:   1,
		Duration:   2,
	},
	Insulted: {
		OpponentHitChance:   20,
		OpponentBlockChance: 20,
		Duration:            3,
	},
}

// ActionString returns the string representation of the action for the condition
func (cd Condition) ActionString() string {
	switch cd {
	case Healthy:
		return "Healthy"
	case Bruised:
		return "Bruise"
	case Disoriented:
		return "Disorientation"
	case Prone:
		return "Knockdown"
	case CriticalHit:
		return "Critical hit"
	case Bleeding:
		return "Bleed"
	case Paralysed:
		return "Paralysis"
	case Insulted:
		return "Insult"
	default:
		return ""
	}
}

// String returns the string representation of the condition
func (cd Condition) String() string {
	switch cd {
	case Healthy:
		return "Healthy"
	case Bruised:
		return "Bruised"
	case Disoriented:
		return "Disoriented"
	case Prone:
		return "Prone"
	case CriticalHit:
		return "Critical Hit"
	case Bleeding:
		return "Bleeding"
	case Paralysed:
		return "Paralysed"
	case Insulted:
		return "Insulted"
	default:
		return ""
	}
}
