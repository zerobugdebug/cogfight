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
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"github.com/sashabaranov/go-openai"

	"github.com/zerobugdebug/cogfight/pkg/attack"
	"github.com/zerobugdebug/cogfight/pkg/modifiers"
	"github.com/zerobugdebug/cogfight/pkg/ui"
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
	TempDamageBonus             float32
	TempComplexityBonus         float32
	TempHitChanceBonus          float32
	TempBlockChanceBonus        float32
	TempSpecialChanceBonus      float32
	Attacks                     []*attack.Attack
	Conditions                  map[modifiers.Condition]int
	CurrentHealth               int
	MaxHealth                   int
}

func (f *Fighter) AddCondition(opponent *Fighter, condition modifiers.Condition) {
	//Add temp bonuses/penalties due to opponent condition
	for modifier, value := range modifiers.DefaultConditionAttributes[condition] {
		switch modifier {
		case modifiers.BlockChance:
			{
				opponent.TempBlockChanceBonus += float32(value)
			}
		case modifiers.HitChance:
			{
				opponent.TempHitChanceBonus += float32(value)
			}
		case modifiers.Damage:
			{
				opponent.TempDamageBonus += float32(value)
			}
		case modifiers.Complexity:
			{
				opponent.TempComplexityBonus += float32(value)
			}
		case modifiers.OpponentHitChance:
			{
				f.TempHitChanceBonus += float32(value)
			}
		case modifiers.OpponentBlockChance:
			{
				f.TempBlockChanceBonus += float32(value)
			}
		}
	}
}

func (f *Fighter) RemoveCondition(opponent *Fighter, condition modifiers.Condition) {
	//Remove temp bonuses/penalties due to opponent condition
	for modifier, value := range modifiers.DefaultConditionAttributes[condition] {
		switch modifier {
		case modifiers.BlockChance:
			{
				opponent.TempBlockChanceBonus -= float32(value)
			}
		case modifiers.HitChance:
			{
				opponent.TempHitChanceBonus -= float32(value)
			}
		case modifiers.Damage:
			{
				opponent.TempDamageBonus -= float32(value)
			}
		case modifiers.Complexity:
			{
				opponent.TempComplexityBonus -= float32(value)
			}
		case modifiers.OpponentHitChance:
			{
				f.TempHitChanceBonus -= float32(value)
			}
		case modifiers.OpponentBlockChance:
			{
				f.TempBlockChanceBonus -= float32(value)
			}
		}
	}
}

func (f *Fighter) ApplyAttack(opponent *Fighter, originalAttack *attack.Attack) {
	modifiedAttack := &attack.Attack{
		Name:           originalAttack.Name,
		Type:           originalAttack.Type,
		Damage:         originalAttack.Damage * (1 + (f.DamageBonus+f.TempDamageBonus)/100),
		Complexity:     originalAttack.Complexity + f.ComplexityBonus + f.TempComplexityBonus,
		HitChance:      originalAttack.HitChance + f.HitChanceBonus + f.TempHitChanceBonus,
		BlockChance:    originalAttack.BlockChance + opponent.BlockChanceBonus + opponent.TempBlockChanceBonus,
		CriticalChance: originalAttack.CriticalChance + f.CriticalChanceBonus,
		SpecialChance:  originalAttack.SpecialChance + f.SpecialChanceBonus + f.TempSpecialChanceBonus,
	}

	sureStrike := 0
	//skipTurn := 0

	//Calculate bonuses/penalties from opponent conditions
	for condition := range opponent.Conditions {
		for modifier, value := range modifiers.DefaultConditionAttributes[condition] {
			switch modifier {
			case modifiers.SureStrike:
				{
					sureStrike = value
				}
			}
		}
	}

	//Calculate bonuses/penalties from attacker conditions
	/* 	for condition := range f.Conditions {
		for modifier, value := range modifiers.DefaultConditionAttributes[condition] {
			switch modifier {
			case modifiers.SkipTurn:
				{
					skipTurn = value
				}
			}
		}

	} */

	var attackDamage float32 = 0
	//Determine if we can attack
	// if skipTurn == 1 {
	// 	fmt.Println("Can't attack, skipping...")
	// } else {
	// Determine the skill of the attacked
	attackComplexity := attack.Clamp(modifiedAttack.Complexity, attack.MinComplexity, attack.MaxComplexity)
	fmt.Printf("Current Complexity: %.1f%%\n", attackComplexity)
	if 100*rand.Float32() > attackComplexity {
		fmt.Println("Attack performed flawlessly!")
		// Determine the attack hit chance
		attackHitChance := attack.Clamp(modifiedAttack.HitChance, attack.MinHitChance, attack.MaxHitChance)
		fmt.Printf("Current Hit Chance: %.1f%%\n", attackHitChance)
		if 100*rand.Float32() < attackHitChance || sureStrike == 1 {
			fmt.Println("Successfull hit!")
			attackBlockChance := attack.Clamp(modifiedAttack.BlockChance, attack.MinBlockChance, attack.MaxBlockChance)
			fmt.Printf("Current Block Chance: %.1f%%\n", attackBlockChance)
			if 100*rand.Float32() > attackBlockChance || sureStrike == 1 {
				fmt.Println("Attack not blocked!")
				attackDamage = attack.Clamp(modifiedAttack.Damage, attack.MinDamage, attack.MaxDamage)
				attackSpecialChance := attack.Clamp(modifiedAttack.SpecialChance, attack.MinSpecialChance, attack.MaxSpecialChance)
				fmt.Printf("Current Special Chance: %.1f%%\n", attackSpecialChance)
				fmt.Printf("Current Special: %s\n", modifiedAttack.Type.Special().ActionString())
				if 100*rand.Float32() < attackSpecialChance {
					fmt.Println("Success! Opponent got " + modifiedAttack.Type.Special().String())
					//if modifiedAttack.Type.Special() == modifiers.CriticalHit {
					//	attackDamage = attack.Clamp(attackDamage*float32(modifiers.DefaultConditionAttributes[modifiedAttack.Type.Special()][modifiers.DamageMult]), attack.MinDamage, attack.MaxDamage)
					//} else {
					_, conditionExist := opponent.Conditions[modifiedAttack.Type.Special()]
					if !conditionExist {
						f.AddCondition(opponent, modifiedAttack.Type.Special())
					}
					opponent.Conditions[modifiedAttack.Type.Special()] = modifiers.DefaultConditionAttributes[modifiedAttack.Type.Special()][modifiers.Duration]

					//}
					/* 						switch modifiedAttack.Type.Special() {
					   						case modifiers.Bleeding:
					   							{
					   								opponent.Conditions[modifiers.Bleeding] = modifiers.DefaultConditionAttributes[modifiers.Bleeding][modifiers.Duration]
					   							}
					   						}
					*/ //attackDamage = clamp(attackDamage*2, MinDamage, MaxDamage)
				} else {
					fmt.Println("Special failed!")
				}
				//fmt.Printf("Damage dealt: %s%.1f%s\n", clrDamage, attackDamage, clrReset)
				//defender.CurrentHealth -= int(attackDamage)
				//fmt.Printf("%s%s takes %d damage! (%d/%d)%s\n", clrBadMessage, defender.Name, int(attackDamage), defender.CurrentHealth, defender.MaxHealth, clrReset)
			} else {
				fmt.Println("Attack blocked!")
			}
		} else {
			fmt.Println("Missed!")
		}
	} else {
		fmt.Println(f.Name, "failed to execute attack!")
	}
	//}

	//Process conditions and specials
	//Calculate effect from opponent conditions
	for condition := range opponent.Conditions {
		for modifier, value := range modifiers.DefaultConditionAttributes[condition] {
			switch modifier {
			case modifiers.DamageMult:
				{
					attackDamage = attackDamage * float32(value)
				}
			}
		}
	}
	if attackDamage > 0 {
		fmt.Printf("Damage dealt: %.1f\n", attackDamage)
		opponent.CurrentHealth -= int(attackDamage)
		fmt.Printf("%s takes %d damage! (%d/%d)\n", opponent.Name, int(attackDamage), opponent.CurrentHealth, opponent.MaxHealth)
	}

	//Calculate effect from attacker conditions
	/* 	for condition := range f.Conditions {
		for modifier, value := range modifiers.DefaultConditionAttributes[condition] {
			switch modifier {
			case modifiers.HPPerTurn:
				{
					f.CurrentHealth += int(value)
					if int(value) < 0 {
						fmt.Printf("%s takes %d damage! (%d/%d) due to %s\n", f.Name, -int(value), f.CurrentHealth, f.MaxHealth, condition.String())
					}
				}
			}
		}
		f.Conditions[condition] -= 1
		if f.Conditions[condition] < 1 {
			delete(f.Conditions, condition)
		}

	} */

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
	fmt.Println("║ Attacks:                                                 ║")
	for i, attack := range f.Attacks {
		fmt.Printf("║ %d. %-53s ║\n", i+1, attack.Name)
	}
	fmt.Println(spacer)
	fmt.Printf("║ Current Health: %-40d ║\n", f.CurrentHealth)
	fmt.Printf("║ Max Health: %-44d ║\n", f.MaxHealth)
	fmt.Println(spacer)
	fmt.Println(bottomBorder)
}

func DisplayFighters(f1, f2 *Fighter) {
	//boxWidth := 50
	numSpacesBetweenFighters := 10
	spaceBetweenFighters := strings.Repeat(" ", numSpacesBetweenFighters)
	var scaleRange float32 = 8.00
	scaleSize := 16

	/*
		blue := color.New(color.FgBlue).SprintFunc()
		red := color.New(color.FgRed).SprintFunc()

		topBorder := "╔" + strings.Repeat("═", boxWidth-2) + "╗"
		bottomBorder := "╚" + strings.Repeat("═", boxWidth-2) + "╝"
		spacer := "║" + strings.Repeat(" ", boxWidth-2) + "║"
		spaceBetweenFighters := strings.Repeat(" ", numSpacesBetweenFighters)

		leftTopBorder := blue(topBorder)
		rightTopBorder := red(topBorder)
		leftBottomBorder := blue(bottomBorder)
		rightBottomBorder := red(bottomBorder)
		leftSpacer := blue(spacer)
		rightSpacer := red(spacer)
		leftEdge := blue("║")
		rightEdge := red("║")
	*/
	/* 	valueFormat
	   	leftFormat := blue("║") + " %-" + strconv.Itoa(halfBoxWidth-4) + "s: %-" + strconv.Itoa(halfBoxWidth-16) + "v " + blue("║") + SpaceBetweenFighters
	   	rightFormat := "║ %-" + strconv.Itoa(halfBoxWidth-4) + "s: %-" + strconv.Itoa(halfBoxWidth-16) + "v ║"
	   	headerFormat := "║ %-" + strconv.Itoa(halfBoxWidth-4) + "s: %-" + strconv.Itoa(halfBoxWidth-16) + "s ║"
	*/
	// ... (previous code)
	/*
		for i, param := range parameters {
			fmt.Printf(leftFormat+rightFormat+"\n", param, values1[i], param, values2[i])
		}
	*/
	blue := color.New(color.BgBlue).SprintFunc()
	//bluefg := color.New(color.FgBlue).SprintFunc()
	red := color.New(color.BgRed).SprintFunc()
	//hired := color.New(color.BgHiRed).SprintFunc()
	hiblue := color.New(color.BgHiBlue).SprintFunc()
	hiblack := color.New(color.BgHiBlack, color.Faint).SprintFunc()
	green := color.New(color.BgGreen).SprintFunc()
	higreen := color.New(color.BgHiGreen).SprintFunc()
	//magenta := color.New(color.BgMagenta).SprintFunc()
	//himagenta := color.New(color.BgHiMagenta).SprintFunc()

	hiredfg := color.New(color.FgHiRed).SprintFunc()
	higreenfg := color.New(color.FgHiGreen).SprintFunc()

	conditionsText := []string{}
	textLeft := []string{}
	var value float32

	textLeft = append(textLeft, "Name: "+f1.Name)
	textLeft = append(textLeft, fmt.Sprintf("Height: %d", f1.Height))
	textLeft = append(textLeft, fmt.Sprintf("Weight: %d", f1.Weight))
	textLeft = append(textLeft, fmt.Sprintf("Age: %d", f1.Age))
	conditionsText = []string{}
	for condition, duration := range f1.Conditions {
		conditionsText = append(conditionsText, fmt.Sprintf("%s[%d]", condition.String(), duration))
	}
	textLeft = append(textLeft, fmt.Sprintf("Conditions: %s", strings.Join(conditionsText, ", ")))
	textLeft = append(textLeft, "")
	textLeft = append(textLeft, fmt.Sprintf("%12s %5.2f %v %5.2f %-12s", "Agility", scaleRange-f1.AgilityStrengthBalance, ui.ScalePrint(-f1.AgilityStrengthBalance, -scaleRange, scaleRange, higreen, hiblue, scaleSize), scaleRange+f1.AgilityStrengthBalance, "Strength"))
	textLeft = append(textLeft, fmt.Sprintf("%12s %5.2f %v %5.2f %-12s", "Burst", scaleRange-f1.BurstEnduranceBalance, ui.ScalePrint(-f1.BurstEnduranceBalance, -scaleRange, scaleRange, higreen, hiblue, scaleSize), scaleRange+f1.BurstEnduranceBalance, "Endurance"))
	textLeft = append(textLeft, fmt.Sprintf("%12s %5.2f %v %5.2f %-12s", "Defense", scaleRange-f1.DefenseOffenseBalance, ui.ScalePrint(-f1.DefenseOffenseBalance, -scaleRange, scaleRange, higreen, hiblue, scaleSize), scaleRange+f1.DefenseOffenseBalance, "Offense"))
	textLeft = append(textLeft, fmt.Sprintf("%12s %5.2f %v %5.2f %-12s", "Speed", scaleRange-f1.SpeedControlBalance, ui.ScalePrint(-f1.SpeedControlBalance, -scaleRange, scaleRange, higreen, hiblue, scaleSize), scaleRange+f1.SpeedControlBalance, "Control"))
	textLeft = append(textLeft, fmt.Sprintf("%12s %5.2f %v %5.2f %-12s", "Intelligence", scaleRange-f1.IntelligenceInstinctBalance, ui.ScalePrint(-f1.IntelligenceInstinctBalance, -scaleRange, scaleRange, higreen, hiblue, scaleSize), scaleRange+f1.IntelligenceInstinctBalance, "Instinct"))
	textLeft = append(textLeft, "")
	value = f1.DamageBonus + f1.TempDamageBonus
	textLeft = append(textLeft, fmt.Sprintf("%20s %7s%% %v", "Damage Bonus", ui.ColorModifiedValue(f1.DamageBonus, value, higreenfg, hiredfg), ui.DoubleScalePrint(value, -100, 0, 100, red, green, hiblack, scaleSize)))
	value = f1.ComplexityBonus + f1.TempComplexityBonus
	textLeft = append(textLeft, fmt.Sprintf("%20s %7s%% %v", "Complexity Bonus", ui.ColorModifiedValue(f1.ComplexityBonus, value, higreenfg, hiredfg), ui.DoubleScalePrint(value, -100, 0, 100, red, green, hiblack, scaleSize)))
	value = f1.HitChanceBonus + f1.TempHitChanceBonus
	textLeft = append(textLeft, fmt.Sprintf("%20s %7s%% %v", "Hit Chance Bonus", ui.ColorModifiedValue(f1.HitChanceBonus, value, higreenfg, hiredfg), ui.DoubleScalePrint(value, -100, 0, 100, red, green, hiblack, scaleSize)))
	value = f1.BlockChanceBonus + f1.TempBlockChanceBonus
	textLeft = append(textLeft, fmt.Sprintf("%20s %7s%% %v", "Block Chance Bonus", ui.ColorModifiedValue(f1.BlockChanceBonus, value, higreenfg, hiredfg), ui.DoubleScalePrint(value, -100, 0, 100, red, green, hiblack, scaleSize)))
	value = f1.SpecialChanceBonus + f1.TempSpecialChanceBonus
	textLeft = append(textLeft, fmt.Sprintf("%20s %7s%% %v", "Special Chance Bonus", ui.ColorModifiedValue(f1.SpecialChanceBonus, value, higreenfg, hiredfg), ui.DoubleScalePrint(value, -100, 0, 100, red, green, hiblack, scaleSize)))
	//textLeft = append(textLeft, fmt.Sprintf("%20s %7.2f%% %v", "Complexity Bonus", f1.ComplexityBonus, ui.DoubleScalePrint(-f1.ComplexityBonus, -100, 0, 100, red, green, hiblack, scaleSize)))
	//textLeft = append(textLeft, fmt.Sprintf("%20s %7.2f%% %v", "Hit Chance Bonus", f1.HitChanceBonus, ui.DoubleScalePrint(f1.HitChanceBonus, -100, 0, 100, red, green, hiblack, scaleSize)))
	//textLeft = append(textLeft, fmt.Sprintf("%20s %7.2f%% %v", "Block Chance Bonus", f1.BlockChanceBonus, ui.DoubleScalePrint(f1.BlockChanceBonus, -100, 0, 100, red, green, hiblack, scaleSize)))
	//textLeft = append(textLeft, fmt.Sprintf("%20s %7.2f%% %v", "Special Chance Bonus", f1.SpecialChanceBonus, ui.DoubleScalePrint(f1.SpecialChanceBonus, -100, 0, 100, red, green, hiblack, scaleSize)))
	textLeft = append(textLeft, "")
	textLeft = append(textLeft, fmt.Sprintf("%s %d/%d %v", "Health: ", f1.CurrentHealth, f1.MaxHealth, ui.ScalePrint(float32(f1.CurrentHealth), 0, float32(f1.MaxHealth), hiblue, hiblack, scaleSize*2)))
	textLeft = append(textLeft, "")

	textRight := []string{}
	textRight = append(textRight, "Name: "+f2.Name)
	textRight = append(textRight, fmt.Sprintf("Height: %d", f2.Height))
	textRight = append(textRight, fmt.Sprintf("Weight: %d", f2.Weight))
	textRight = append(textRight, fmt.Sprintf("Age: %d", f2.Age))
	conditionsText = []string{}
	for condition, duration := range f2.Conditions {
		conditionsText = append(conditionsText, fmt.Sprintf("%s[%d]", condition.String(), duration))
	}
	textRight = append(textRight, fmt.Sprintf("Conditions: %s", strings.Join(conditionsText, ", ")))
	textRight = append(textRight, "")
	textRight = append(textRight, fmt.Sprintf("%12s %5.2f %v %5.2f %-12s", "Agility", scaleRange-f2.AgilityStrengthBalance, ui.ScalePrint(-f2.AgilityStrengthBalance, -scaleRange, scaleRange, higreen, hiblue, scaleSize), scaleRange+f2.AgilityStrengthBalance, "Strength"))
	textRight = append(textRight, fmt.Sprintf("%12s %5.2f %v %5.2f %-12s", "Burst", scaleRange-f2.BurstEnduranceBalance, ui.ScalePrint(-f2.BurstEnduranceBalance, -scaleRange, scaleRange, higreen, hiblue, scaleSize), scaleRange+f2.BurstEnduranceBalance, "Endurance"))
	textRight = append(textRight, fmt.Sprintf("%12s %5.2f %v %5.2f %-12s", "Defense", scaleRange-f2.DefenseOffenseBalance, ui.ScalePrint(-f2.DefenseOffenseBalance, -scaleRange, scaleRange, higreen, hiblue, scaleSize), scaleRange+f2.DefenseOffenseBalance, "Offense"))
	textRight = append(textRight, fmt.Sprintf("%12s %5.2f %v %5.2f %-12s", "Speed", scaleRange-f2.SpeedControlBalance, ui.ScalePrint(-f2.SpeedControlBalance, -scaleRange, scaleRange, higreen, hiblue, scaleSize), scaleRange+f2.SpeedControlBalance, "Control"))
	textRight = append(textRight, fmt.Sprintf("%12s %5.2f %v %5.2f %-12s", "Intelligence", scaleRange-f2.IntelligenceInstinctBalance, ui.ScalePrint(-f2.IntelligenceInstinctBalance, -scaleRange, scaleRange, higreen, hiblue, scaleSize), scaleRange+f2.IntelligenceInstinctBalance, "Instinct"))
	textRight = append(textRight, "")
	value = f2.DamageBonus + f2.TempDamageBonus
	textRight = append(textRight, fmt.Sprintf("%20s %7s%% %v", "Damage Bonus", ui.ColorModifiedValue(f2.DamageBonus, value, higreenfg, hiredfg), ui.DoubleScalePrint(value, -100, 0, 100, red, green, hiblack, scaleSize)))
	value = f2.ComplexityBonus + f2.TempComplexityBonus
	textRight = append(textRight, fmt.Sprintf("%20s %7s%% %v", "Complexity Bonus", ui.ColorModifiedValue(f2.ComplexityBonus, value, higreenfg, hiredfg), ui.DoubleScalePrint(value, -100, 0, 100, red, green, hiblack, scaleSize)))
	value = f2.HitChanceBonus + f2.TempHitChanceBonus
	textRight = append(textRight, fmt.Sprintf("%20s %7s%% %v", "Hit Chance Bonus", ui.ColorModifiedValue(f2.HitChanceBonus, value, higreenfg, hiredfg), ui.DoubleScalePrint(value, -100, 0, 100, red, green, hiblack, scaleSize)))
	value = f2.BlockChanceBonus + f2.TempBlockChanceBonus
	textRight = append(textRight, fmt.Sprintf("%20s %7s%% %v", "Block Chance Bonus", ui.ColorModifiedValue(f2.BlockChanceBonus, value, higreenfg, hiredfg), ui.DoubleScalePrint(value, -100, 0, 100, red, green, hiblack, scaleSize)))
	value = f2.SpecialChanceBonus + f2.TempSpecialChanceBonus
	textRight = append(textRight, fmt.Sprintf("%20s %7s%% %v", "Special Chance Bonus", ui.ColorModifiedValue(f2.SpecialChanceBonus, value, higreenfg, hiredfg), ui.DoubleScalePrint(value, -100, 0, 100, red, green, hiblack, scaleSize)))

	// textRight = append(textRight, fmt.Sprintf("%20s %7.2f%% %v", "Damage Bonus", f2.DamageBonus, ui.DoubleScalePrint(f2.DamageBonus, -100, 0, 100, red, green, hiblack, scaleSize)))
	// textRight = append(textRight, fmt.Sprintf("%20s %7.2f%% %v", "Complexity Bonus", f2.ComplexityBonus, ui.DoubleScalePrint(-f2.ComplexityBonus, -100, 0, 100, red, green, hiblack, scaleSize)))
	// textRight = append(textRight, fmt.Sprintf("%20s %7.2f%% %v", "Hit Chance Bonus", f2.HitChanceBonus, ui.DoubleScalePrint(f2.HitChanceBonus, -100, 0, 100, red, green, hiblack, scaleSize)))
	// textRight = append(textRight, fmt.Sprintf("%20s %7.2f%% %v", "Block Chance Bonus", f2.BlockChanceBonus, ui.DoubleScalePrint(f2.BlockChanceBonus, -100, 0, 100, red, green, hiblack, scaleSize)))
	// textRight = append(textRight, fmt.Sprintf("%20s %7.2f%% %v", "Special Chance Bonus", f2.SpecialChanceBonus, ui.DoubleScalePrint(f2.SpecialChanceBonus, -100, 0, 100, red, green, hiblack, scaleSize)))
	textRight = append(textRight, "")
	textRight = append(textRight, fmt.Sprintf("%s %d/%d %v", "Health: ", f2.CurrentHealth, f2.MaxHealth, ui.ScalePrint(float32(f2.CurrentHealth), 0, float32(f2.MaxHealth), hiblue, hiblack, scaleSize*2)))
	textRight = append(textRight, "")

	boxLeft := ui.BoxPrint(20, blue, textLeft)
	boxRight := ui.BoxPrint(20, red, textRight)

	for i, v := range boxLeft {
		fmt.Println(v + spaceBetweenFighters + boxRight[i])
	}

	//	fmt.Println(leftTopBorder + spaceBetweenFighters + rightTopBorder)
	//	fmt.Println(leftSpacer + spaceBetweenFighters + rightSpacer)
	/*	textLeft = " Name: " + f1.Name + strings.Repeat(" ", boxWidth-len(f1.Name)-2-len(" Name: "))
		textRight = " Name: " + f2.Name + strings.Repeat(" ", boxWidth-len(f2.Name)-2-len(" Name: "))
		fmt.Println(leftEdge + textLeft + leftEdge + spaceBetweenFighters + rightEdge + textRight + rightEdge)
		textLeft = " Height: " + strconv.Itoa(f1.Height) + strings.Repeat(" ", boxWidth-len(strconv.Itoa(f1.Height))-2-len(" Height: "))
		textRight = " Height: " + strconv.Itoa(f2.Height) + strings.Repeat(" ", boxWidth-len(strconv.Itoa(f2.Height))-2-len(" Height: "))
		fmt.Println(leftEdge + textLeft + leftEdge + spaceBetweenFighters + rightEdge + textRight + rightEdge)
		textLeft = " Weight: " + strconv.Itoa(f1.Weight) + strings.Repeat(" ", boxWidth-len(strconv.Itoa(f1.Weight))-2-len(" Weight: "))
		textRight = " Weight: " + strconv.Itoa(f2.Weight) + strings.Repeat(" ", boxWidth-len(strconv.Itoa(f2.Weight))-2-len(" Weight: "))
		fmt.Println(leftEdge + textLeft + leftEdge + spaceBetweenFighters + rightEdge + textRight + rightEdge)
		textLeft = " Age: " + strconv.Itoa(f1.Age) + strings.Repeat(" ", boxWidth-len(strconv.Itoa(f1.Age))-2-len(" Age: "))
		textRight = " Age: " + strconv.Itoa(f2.Age) + strings.Repeat(" ", boxWidth-len(strconv.Itoa(f2.Age))-2-len(" Age: "))
		fmt.Println(leftEdge + textLeft + leftEdge + spaceBetweenFighters + rightEdge + textRight + rightEdge)
		textLeft = "Agility " + strconv.FormatFloat(float64(12-f1.AgilityStrengthBalance), 'f', 1, 32) + ui.ScalePrint(f1.AgilityStrengthBalance, -10, 10, color.New(color.BgBlue).SprintFunc(), color.New(color.BgRed).SprintFunc(), 20) + " " + strconv.FormatFloat(float64(8+f1.AgilityStrengthBalance), 'f', 1, 32) + " Strength"
		textRight = ui.ScalePrint(f2.AgilityStrengthBalance, -10, 10, color.New(color.BgBlue).SprintFunc(), color.New(color.BgRed).SprintFunc(), 20)
		fmt.Println(leftEdge + textLeft + leftEdge + spaceBetweenFighters + rightEdge + textRight + rightEdge)*/
	/*	fmt.Println(headerFormat, "Attacks", "Attacks")

		maxAttacks := max(len(f1.Attacks), len(f2.Attacks))
		for i := 0; i < maxAttacks; i++ {
			attack1 := " "
			attack2 := " "
			if i < len(f1.Attacks) {
				attack1 = f1.Attacks[i].Name
			}
			if i < len(f2.Attacks) {
				attack2 = f2.Attacks[i].Name
			}
			fmt.Printf(leftFormat+rightFormat+"\n", strconv.Itoa(i+1)+". "+attack1, "", strconv.Itoa(i+1)+". "+attack2, "")
		}
	*/
	/* 	fmt.Println(leftSpacer + spaceBetweenFighters + rightSpacer)
	   	fmt.Println(leftBottomBorder + spaceBetweenFighters + rightBottomBorder) */

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
			Message: "Enter fighter height (160-200 cm):",
			Help:    "Please enter your fighter height. Taller fighters will favour Strength, Offense and Control, while lower height will give benefits to Agility, Defense and Speed.",
			Default: "180",
		},
		Validate: validateNumber(160, 200),
	}
	qs = append(qs, heightQuestion)

	weightQuestion := &survey.Question{
		Name: "weight",
		Prompt: &survey.Input{
			Message: "Enter fighter weight (60-120 kg):",
			Help:    "Please enter your fighter weight. Heavier fighters tend to have better Strength, Endurance and Control, while lighter fighters rely more on the Agility, Burst and Speed.",
			Default: "90",
		},
		Validate: validateNumber(60, 120),
	}
	qs = append(qs, weightQuestion)

	ageQuestion := &survey.Question{
		Name: "age",
		Prompt: &survey.Input{
			Message: "Enter fighter age (18-60 years):",
			Help:    "Please enter your fighter age. Older fighters tend to have better Intelligence, while younger fighters rely more on the Instinct.",
			Default: "40",
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
			Default: "Balanced",
		},
	}
	qs = append(qs, agilityStrengthBalanceQuestion)

	burstEnduranceBalanceQuestion := &survey.Question{
		Name: "burstEnduranceBalance",
		Prompt: &survey.Select{
			Message: "Choose fighter burst/endurance balance:",
			Options: []string{"Very high Burst, Very low Endurance", "High Burst, Low Endurance", "Balanced", "Low Burst, High Endurance", "Very low Burst, Very high Endurance"},
			Help:    "This parameter determines the balance between Burst and Endurance. Fighters with high Burst will get better chances to hit and special effects, but high Endurance will give bonuses to damage and blocking chance.",
			Default: "Balanced",
		},
	}
	qs = append(qs, burstEnduranceBalanceQuestion)

	defenseOffenseBalanceQuestion := &survey.Question{
		Name: "defenseOffenseBalance",
		Prompt: &survey.Select{
			Message: "Choose fighter defense/offense balance:",
			Options: []string{"Very high Defense, Very low Offense", "High Defense, Low Offense", "Balanced", "Low Defense, High Offense", "Very low Defense, Very high Offense"},
			Help:    "This parameter determines the balance between Defense and Offense. Increasing Defense will improve your chances of blocking attacks, while increasing Offense will help with hitting.",
			Default: "Balanced",
		},
	}
	qs = append(qs, defenseOffenseBalanceQuestion)

	speedControlBalanceQuestion := &survey.Question{
		Name: "speedControlBalance",
		Prompt: &survey.Select{
			Message: "Choose fighter speed/control balance:",
			Options: []string{"Very high Speed, Very low Control", "High Speed, Low Control", "Balanced", "Low Speed, High Control", "Very low Speed, Very high Control"},
			Help:    "This parameter determines the balance between Speed and Control. Increasing Speed will improve your chances of successfully hitting and blocking attacks, while high Control will help with executing more complex attacks and critical hits",
			Default: "Balanced",
		},
	}
	qs = append(qs, speedControlBalanceQuestion)

	intelligenceInstinctBalanceQuestion := &survey.Question{
		Name: "intelligenceInstinctBalance",
		Prompt: &survey.Select{
			Message: "Choose fighter intelligence/instinct balance:",
			Options: []string{"Very high Intelligence, Very low Instinct", "High Intelligence, Low Instinct", "Balanced", "Low Intelligence, High Instinct", "Very low Intelligence, Very high Instinct"},
			Help:    "This parameter determines the balance between Intelligence and Instinct. Increasing Intelligence will help with executing more complex attacks and critical hits, while increasing Instinct will improve your chances of successfully hitting and blocking attacks",
			Default: "Balanced",
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
		fmt.Println("Error during the fighter creation:", err)
		return nil
	}

	fmt.Println(answers)

	attackType := attack.AttackType(0)
	attackTypePromptOptions := []string{}

	for attackType.String() != "" {
		attackTypePromptOptions = append(attackTypePromptOptions, attackType.String())
		//+": "+attackType.Hint())
		attackType++
	}

	attackTypePrompt := &survey.Select{
		Message:  "Select an attack type:",
		Options:  attackTypePromptOptions,
		PageSize: attack.MaxAttackTypes,
		Help:     "Punch: Closed fist attacks, high damage, low complexity, high hit chance, high block chance\nSlap: Open fist or back hand attacks, very low damage, low complexity, high hit chance, high block chance\nKick: Leg attacks, high damage, average complexity, high hit chance, high block chance\nKnee strike: Attacks with a knee, very high damage, average complexity, high hit chance, average block chance\nElbow strike: Attacks with an elbow, very high damage, low complexity, high hit chance, high block chance\nThrow: Attacks to knockdown opponent, average damage, average complexity, average hit chance, average block chance, can knockdown opponent\nLock: Grapple attacks to block joint movement, very low damage, high complexity, low hit chance, low block chance, decrease opponent's hit and block chances\nChoke: Grapple attacks to block airways, low damage, high complexity, low hit chance, low block chance, decrease opponent's damage and increase complexity\nCustom: Custom free text attack",
	}

	// Create the fighter object
	fighter := &Fighter{
		Name:                        answers.Name,
		Height:                      answers.Height,
		Weight:                      answers.Weight,
		Age:                         answers.Age,
		AgilityStrengthBalance:      float32(answers.AgilityStrengthBalance) + (float32(answers.Weight)-90)/15 + (float32(answers.Height)-180)/10 - 2,
		BurstEnduranceBalance:       float32(answers.BurstEnduranceBalance) + (float32(answers.Weight)-90)/15 - 2,
		DefenseOffenseBalance:       float32(answers.DefenseOffenseBalance) + (float32(answers.Height)-180)/10 - 2,
		SpeedControlBalance:         float32(answers.SpeedControlBalance) + (float32(answers.Weight)-90)/15 + (float32(answers.Height)-180)/10 - 2,
		IntelligenceInstinctBalance: float32(answers.IntelligenceInstinctBalance) - (float32(answers.Age)-39)/10 - 2,
		CurrentHealth:               250 + (answers.Weight - 90),
		MaxHealth:                   250 + (answers.Weight - 90),
		Conditions:                  make(map[modifiers.Condition]int),
	}

	fighter.DamageBonus = (fighter.AgilityStrengthBalance + fighter.BurstEnduranceBalance) * 10
	fighter.ComplexityBonus = (fighter.AgilityStrengthBalance - fighter.SpeedControlBalance + fighter.IntelligenceInstinctBalance) * 5
	fighter.HitChanceBonus = (-fighter.AgilityStrengthBalance - fighter.BurstEnduranceBalance + fighter.DefenseOffenseBalance - fighter.SpeedControlBalance + fighter.IntelligenceInstinctBalance) * 3
	fighter.BlockChanceBonus = (-fighter.AgilityStrengthBalance + fighter.BurstEnduranceBalance - fighter.DefenseOffenseBalance - fighter.SpeedControlBalance + fighter.IntelligenceInstinctBalance) * 3
	fighter.CriticalChanceBonus = (fighter.SpeedControlBalance - fighter.IntelligenceInstinctBalance) * 10
	fighter.SpecialChanceBonus = (fighter.AgilityStrengthBalance - fighter.BurstEnduranceBalance) * 10

	defaultAttacks := attack.NewDefaultAttacks()
	attacks := []*attack.Attack{}

	i := 0
	for i < numAttacks {
		fmt.Printf("Attack %d from %d\n", i+1, numAttacks)
		attackTypeSelected := 0
		// Ask for attack type
		err := survey.AskOne(attackTypePrompt, &attackTypeSelected, survey.WithValidator(survey.Required))
		if err != nil {
			fmt.Println("Error during the attack type selection:", err)
			break
		}
		attackType = attack.AttackType(attackTypeSelected)
		fmt.Println("defaultAttacks.GetAttacksByType(attackType)=", defaultAttacks.GetAttacksByType(attackType))

		// If non-custom type, ask for specific attack
		if attackType != attack.Custom {
			attackNamePromptOptions := []string{}
			for _, value := range defaultAttacks.GetAttacksByType(attackType) {
				//				damage:=clamp(value.Damage+fighter.DamageBonus,)
				//attackDescription := fmt.Sprintf("%s [DMG: %5.2f, CMP: %5.2f, HIT: %5.2f, BLK: %5.2f, SPC: %5.2f]", value.Name, value.Damage+fighter.DamageBonus, value.Complexity+fighter.ComplexityBonus, value.HitChance+fighter.HitChanceBonus, value.BlockChance+fighter.BlockChanceBonus, value.SpecialChance+fighter.SpecialChanceBonus)
				attackNamePromptOptions = append(attackNamePromptOptions, value.Name)
			}
			attackNamePromptOptions = append(attackNamePromptOptions, "<-Back")

			attackNamePrompt := &survey.Select{
				Message:  "Select an attack:",
				Options:  attackNamePromptOptions,
				PageSize: len(attackNamePromptOptions),
				Description: func(value string, index int) string {
					if value != "<-Back" {
						attack := defaultAttacks.GetAttackByName(value)
						return fmt.Sprintf("[DMG: %5.2f, CMP: %5.2f, HIT: %5.2f, BLK: %5.2f, SPC: %5.2f]", attack.Damage+fighter.DamageBonus, attack.Complexity+fighter.ComplexityBonus, attack.HitChance+fighter.HitChanceBonus, attack.BlockChance+fighter.BlockChanceBonus, attack.SpecialChance+fighter.SpecialChanceBonus)
					}
					return ""
				},
			}
			attackName := ""
			err = survey.AskOne(attackNamePrompt, &attackName, survey.WithValidator(survey.Required))
			if err != nil {
				fmt.Println("Error during the attack selection:", err)
				break
			}
			if attackName != "<-Back" {
				// Add attack to attacks array
				attacks = append(attacks, defaultAttacks.GetAttackByName(attackName))
				fmt.Printf("attacks= %v\n", attacks)
				i++
			}
			continue
		}

		// Ask for the description of the custom attack
		customAttackName := ""
		customAttackPrompt := &survey.Input{
			Message: "Enter the description for the custom attack:",
		}
		err = survey.AskOne(customAttackPrompt, &customAttackName)
		fmt.Println("customAttackName=" + customAttackName)
		if err != nil {
			fmt.Println("Error during the custom attack creation:", err)
			break
		}

		// Validate the attack name and get the attack parameters using OpenAI API
		validAttack, err := validateAttackName(customAttackName)
		if err != nil {
			fmt.Printf("Error validating attack name: %s\n", err)
			continue
		}

		if validAttack {
			fmt.Println("Valid attack")
			customAttackType, err := getOpenAIResponse("COG_TYPE_ATTACK_PROMPT", customAttackName)
			if err != nil {
				fmt.Printf("Error getting data for COG_COMPLEXITY_ATTACK_PROMPT: %s\n", err)
				continue
			}
			fmt.Println("customAttackType=", customAttackType.(string))

			complexityValue, err := getOpenAIResponse("COG_COMPLEXITY_ATTACK_PROMPT", customAttackName)
			if err != nil {
				fmt.Printf("Error getting data for COG_COMPLEXITY_ATTACK_PROMPT: %s\n", err)
				continue
			}
			fmt.Println("complexityValue=", complexityValue.(int))
		}
	}
	fmt.Printf("attacks_outside= %v\n", attacks)

	// Create the fighter attacks
	fighter.Attacks = append(fighter.Attacks, attacks...)
	fmt.Printf("fighter.Attacks= %v\n", fighter.Attacks)

	fmt.Printf("\n%s has been created!\n", fighter.Name)
	fighter.DisplayFighter()

	return fighter
}

// GenerateComputerFighter generates a computer-controlled fighter
func GenerateComputerFighter(playerFighter *Fighter) *Fighter {
	rand.Seed(time.Now().UnixNano())

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

	answers.Height = rand.Intn(41) + 160 // Height between 160 and 200 cm
	answers.Weight = rand.Intn(61) + 60  // Weight between 60 and 120 kg
	answers.Age = rand.Intn(43) + 18     // Age between 18 and 60 years

	// Generate random values for the computer fighter's attributes
	computerFighter := &Fighter{
		Name:                        "Computer Fighter",
		Height:                      answers.Height,
		Weight:                      answers.Weight,
		Age:                         answers.Age,
		AgilityStrengthBalance:      float32(answers.AgilityStrengthBalance) + (float32(answers.Weight)-90)/15 + (float32(answers.Height)-180)/10 - 2,
		BurstEnduranceBalance:       float32(answers.BurstEnduranceBalance) + (float32(answers.Weight)-90)/15 - 2,
		DefenseOffenseBalance:       float32(answers.DefenseOffenseBalance) + (float32(answers.Height)-180)/10 - 2,
		SpeedControlBalance:         float32(answers.SpeedControlBalance) + (float32(answers.Weight)-90)/15 + (float32(answers.Height)-180)/10 - 2,
		IntelligenceInstinctBalance: float32(answers.IntelligenceInstinctBalance) - (float32(answers.Age)-39)/10 - 2,
		Attacks:                     []*attack.Attack{},
		CurrentHealth:               250 + (answers.Weight - 90),
		MaxHealth:                   250 + (answers.Weight - 90),
		Conditions:                  make(map[modifiers.Condition]int),
	}

	computerFighter.DamageBonus = (computerFighter.AgilityStrengthBalance + computerFighter.BurstEnduranceBalance) * 10
	computerFighter.ComplexityBonus = (computerFighter.AgilityStrengthBalance - computerFighter.SpeedControlBalance + computerFighter.IntelligenceInstinctBalance) * 5
	computerFighter.HitChanceBonus = (-computerFighter.AgilityStrengthBalance - computerFighter.BurstEnduranceBalance + computerFighter.DefenseOffenseBalance - computerFighter.SpeedControlBalance + computerFighter.IntelligenceInstinctBalance) * 3
	computerFighter.BlockChanceBonus = (-computerFighter.AgilityStrengthBalance + computerFighter.BurstEnduranceBalance - computerFighter.DefenseOffenseBalance - computerFighter.SpeedControlBalance + computerFighter.IntelligenceInstinctBalance) * 3
	computerFighter.CriticalChanceBonus = (computerFighter.SpeedControlBalance - computerFighter.IntelligenceInstinctBalance) * 10
	computerFighter.SpecialChanceBonus = (computerFighter.AgilityStrengthBalance - computerFighter.BurstEnduranceBalance) * 10

	defaultAttacks := attack.NewDefaultAttacks()
	for range playerFighter.Attacks {
		attacksList := defaultAttacks.GetAttacksByType(attack.AttackType(rand.Intn(attack.MaxAttackTypes - 1)))
		fmt.Println("attacksList=", attacksList)
		fmt.Println("len(attacksList)=", len(attacksList))
		computerAttack := attacksList[rand.Intn(len(attacksList))]
		fmt.Println("computerAttack=", computerAttack)
		computerFighter.Attacks = append(computerFighter.Attacks, computerAttack)
	}

	fmt.Printf("\n%s has been generated!\n", computerFighter.Name)
	computerFighter.DisplayFighter()
	return computerFighter
}

// validateAttackName validates the given attack name using OpenAI API and returns the attack parameters
func validateAttackName(attackName string) (bool, error) {

	attackValidation, err := getOpenAIResponse("COG_VALIDATION_ATTACK_PROMPT", attackName)
	if err != nil {
		return false, fmt.Errorf("error sending OpenAI API request: %s", err)
	}
	/* 	apiKey := os.Getenv("OPENAI_API_KEY")
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
	   		return false, fmt.Errorf("error sending OpenAI API request: %s", err)
	   	}

	   	fmt.Println(response.Choices[0].Message.Content)
	   	// Parse the response and extract integer answer
	   	reply := response.Choices[0].Message.Content

	   	// Parse the response to confirm if attack is valid
	   	//reply := response.Choices[0].Message.Content
	   	client = nil
	*/fmt.Println("attackValidation=", attackValidation)
	reply := attackValidation.(string)

	if strings.Contains(reply, "Invalid") {
		return false, errors.New("Not a valid attack")
		// } else if strings.Contains(reply, "Multiple") {
		// 	return false, errors.New("Attack is valid, but not a single attack")
	} else if strings.Contains(reply, "Impossible") || strings.Contains(reply, "Valid") || strings.Contains(reply, "Multiple") {
		return true, nil
	}

	return false, fmt.Errorf("Unknown response from OpenAI API: %s", reply)
}

// Get integer answer from OpenAI API
func getOpenAIResponse(promptEnvVariable string, promptData string) (interface{}, error) {
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
			Model:     openai.GPT3Dot5Turbo,
			MaxTokens: 1000,
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

	fmt.Println("response.Choices[0].Message.Content=", response.Choices[0].Message.Content)
	// Parse the response and extract integer answer
	reply := response.Choices[0].Message.Content
	re := regexp.MustCompile(`\[\[(\d+)\]\]`)
	matchInt := re.FindStringSubmatch(reply)
	fmt.Println("matchInt=", matchInt)
	if len(matchInt) > 1 {
		fmt.Println("Number:", matchInt[1])
		return strconv.Atoi(matchInt[1])
	}
	re = regexp.MustCompile(`\[\[(\w+)\]\]`)
	matchString := re.FindStringSubmatch(reply)
	fmt.Println("matchString=", matchString)
	if len(matchString) > 1 {
		fmt.Println("String:", matchString[1])
		return matchString[1], nil
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
