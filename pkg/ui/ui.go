package ui

import (
	"fmt"
	"math"
	"regexp"
	"strings"
	"sync"
	"time"
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

func ColorModifiedValue(value, delta float32, format string, colorMore func(a ...interface{}) string, colorLess func(a ...interface{}) string) string {
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
