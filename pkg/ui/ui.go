package ui

import (
	"fmt"
	"math"
	"regexp"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/fatih/color"
)

type Alignment int

const (
	Left Alignment = iota
	Center
	Right
)

type colorDelims struct {
	Start, End string
	Color      *color.Color
	Remove     bool
}

// Const map alternative
func colorDelimMap() map[string]colorDelims {
	return map[string]colorDelims{
		"cyan":   {"[", "]", color.New(color.FgCyan), false},
		"yellow": {"{", "}", color.New(color.FgYellow), true},
		"green":  {"\"", "\"", color.New(color.FgGreen), true},
	}
}

// Const *struct alternative
func defaultColor() *color.Color {
	return color.New(color.FgWhite)
}

func ColorizeChunk(chunk string, stateStack []string) (string, []string) {
	coloredChunk := ""
	colorMap := colorDelimMap()
	for i := 0; i < len(chunk); i++ {
		char := string(chunk[i])
		currentState := stateStack[len(stateStack)-1]
		if currentState == "default" {
			for colorKey, delimiters := range colorMap {
				if char == delimiters.Start {
					if !delimiters.Remove {
						coloredChunk += delimiters.Color.Sprint(char)
						//delimiters.Color.Print(char)
					}
					//Push to stack
					stateStack = append(stateStack, colorKey)
					currentState = colorKey
					break
				}
			}
			if currentState == "default" {
				coloredChunk += defaultColor().Sprint(char)
			}
		} else {
			if char == colorMap[currentState].End {
				if !colorMap[currentState].Remove {
					coloredChunk += colorMap[currentState].Color.Sprint(char)
				}
				//Pop from stack
				stateStack = stateStack[:len(stateStack)-1]
				if len(stateStack) == 0 {
					//color.New(color.FgRed).Println("Error: Unmatched delimiter")
					//Push to stack
					stateStack = append(stateStack, "default")
				}
			} else {
				coloredChunk += colorMap[currentState].Color.Sprint(char)
			}
		}
	}

	/* 	// Error Handling
	   	if len(stateStack) > 1 {
	   		color.New(color.FgRed).Println("Error: Unmatched delimiter")
	   		stateStack = []string{"normal"}
	   	} */

	return coloredChunk, stateStack
}

func ScalePrint(value, min, max float64, colorLeft func(a ...interface{}) string, colorRight func(a ...interface{}) string, length int) string {
	normalizedValue := (value - min) / (max - min)
	position := int(math.Round(float64(normalizedValue * float64(length))))

	leftString := colorLeft(strings.Repeat(" ", position))
	//middleString := "|"
	rightString := colorRight(strings.Repeat(" ", length-position))

	return leftString + rightString
}

func DoubleScalePrint(value, min, center, max float64, colorLeft func(a ...interface{}) string, colorRight func(a ...interface{}) string, colorBack func(a ...interface{}) string, length int) string {

	if value >= center {
		return ScalePrint(center, min, center, colorBack, colorLeft, length/2) + colorRight("\x1B[30m│\x1B[0m") + ScalePrint(value, center, max, colorRight, colorBack, length/2)
	} else {
		return ScalePrint(value, min, center, colorBack, colorLeft, length/2) + colorLeft("\x1B[30m│\x1B[0m") + ScalePrint(center, center, max, colorRight, colorBack, length/2)
	}
}

func AlignText(text string, size int, optionalArgs ...interface{}) string {
	fillChar := ' '
	alignment := Center

	for _, arg := range optionalArgs {
		switch v := arg.(type) {
		case rune:
			fillChar = v
		case Alignment:
			alignment = v
		}
	}

	ansiRegex := regexp.MustCompile(`\x1B\[[0-?]*[ -/]*[@-~]`)
	textLength := utf8.RuneCountInString(ansiRegex.ReplaceAllString(text, ""))
	if textLength >= size {
		return text
	}

	padding := size - textLength
	leftPadding := 0
	rightPadding := 0

	switch alignment {
	case Left:
		rightPadding = padding
	case Center:
		leftPadding = padding / 2
		rightPadding = padding - leftPadding
	case Right:
		leftPadding = padding
	}

	return strings.Repeat(string(fillChar), leftPadding) + text + strings.Repeat(string(fillChar), rightPadding)
}

func BoxPrint(minWidth int, colorFunc func(a ...interface{}) string, lines []string) []string {
	maxWidth := minWidth

	for _, line := range lines {
		ansiRegex := regexp.MustCompile(`\x1B\[[0-?]*[ -/]*[@-~]`)
		pureLine := ansiRegex.ReplaceAllString(line, "")
		if utf8.RuneCountInString(pureLine) > maxWidth {
			maxWidth = len(pureLine)
			//fmt.Println("line=", pureLine)
			//fmt.Println("maxWidth=", maxWidth)
		}
	}

	//topBorder := colorFunc(strings.Repeat("═", maxWidth+2))
	//bottomBorder := colorFunc(strings.Repeat("═", maxWidth+2))

	topBorder := colorFunc("╔" + strings.Repeat("═", maxWidth+2) + "╗")
	bottomBorder := colorFunc("╚" + strings.Repeat("═", maxWidth+2) + "╝")
	spacer := colorFunc("║")

	box := []string{topBorder}
	for _, line := range lines {
		box = append(box, spacer+" "+AlignText(line, maxWidth)+" "+spacer)
	}
	box = append(box, bottomBorder)

	return box
}

func ColorModifiedValue(value, delta float64, format string, colorMore func(a ...interface{}) string, colorLess func(a ...interface{}) string) string {
	if delta < 0 {
		return colorLess(fmt.Sprintf(format, value))
	} else {
		if delta > 0 {
			return colorMore(fmt.Sprintf(format, value))
		} else {

			return fmt.Sprintf(format, value)
		}
	}
}

func RotatingPipe(stopChan chan bool, wg *sync.WaitGroup) {
	defer wg.Done()
	pipeChars := []string{"|", "/", "-", "\\"}
	i := 0

	for {
		select {
		case <-stopChan:
			//fmt.Println("\nRotating pipe stopped.")
			return
		default:
			fmt.Printf("\r%s", pipeChars[i])
			i = (i + 1) % len(pipeChars)
			time.Sleep(100 * time.Millisecond)
		}
	}
}
