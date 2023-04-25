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

// ConditionAttributes maps fighter conditions to their respective attributes
var ConditionAttributes = map[Condition]map[string]int{
	Healthy: {},
	Disoriented: {
		"hitChance":      -20,
		"blockChance":    -20,
		"takedownChance": 20,
	},
	Choked: {
		"hitChance":      -30,
		"blockChance":    -30,
		"takedownChance": 0,
	},
}
