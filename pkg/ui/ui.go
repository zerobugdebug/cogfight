package ui

import (
	"fmt"
	"math"
	"regexp"
	"strings"
)

type Alignment int

const (
	Left Alignment = iota
	Center
	Right
)

func ScalePrint(value, min, max float32, colorLeft func(a ...interface{}) string, colorRight func(a ...interface{}) string, length int) string {
	normalizedValue := (value - min) / (max - min)
	position := int(math.Round(float64(normalizedValue * float32(length))))

	leftString := colorLeft(strings.Repeat(" ", position))
	//middleString := "|"
	rightString := colorRight(strings.Repeat(" ", length-position))

	return leftString + rightString
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
	textLength := len(ansiRegex.ReplaceAllString(text, ""))
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
		if len(pureLine) > maxWidth {
			maxWidth = len(pureLine)
			fmt.Println("line=", pureLine)
			fmt.Println("maxWidth=", maxWidth)
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
