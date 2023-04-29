package ui

import (
	"fmt"
	"math"
	"regexp"
	"strings"
	"unicode/utf8"
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

func DoubleScalePrint(value, min, center, max float32, colorLeft func(a ...interface{}) string, colorRight func(a ...interface{}) string, colorBack func(a ...interface{}) string, length int) string {

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

func ColorModifiedValue(old, new float32, colorMore func(a ...interface{}) string, colorLess func(a ...interface{}) string) string {
	if int((old-new)*100) > 0 {
		return colorLess(fmt.Sprintf("%.2f", new))
	} else {
		if int((old-new)*100) < 0 {
			return colorMore(fmt.Sprintf("%.2f", new))
		} else {

			return fmt.Sprintf("%.2f", old)
		}
	}
}
