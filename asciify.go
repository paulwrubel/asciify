package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Letter struct {
	character rune
	rowCount  int
	rowWidth  int
	rows      []string
}

func main() {
	globalRowCount := 6

	fontFileName := flag.String("file", "font.csv", "(required) Path to font file location")
	lettersPerRow := flag.Int("width", -1, "(optional) Letters per line")
	splitCharacter := flag.String("split", "", "(optional) Character to split lines on")

	flag.Parse()

	rawInput := flag.Arg(0)

	fontFile, err := os.Open(*fontFileName)
	if err != nil {
		fmt.Printf("Error opening font file: %s\n", err.Error())
		return
	}
	fontReader := csv.NewReader(fontFile)
	fontData, err := fontReader.ReadAll()
	if err != nil {
		fmt.Printf("Error reading font file: %s\n", err.Error())
		return
	}
	letterMap := map[rune]*Letter{}
	for _, row := range fontData[1:] {
		letterRune := rune(row[0][0])
		rowNum, err := strconv.Atoi(row[1])
		if err != nil {
			fmt.Printf("Error parsing int from string: %s\n", err.Error())
			return
		}
		rowString := row[2]
		rowWidth := len(rowString)
		if err != nil {
			fmt.Printf("Error parsing int from string: %s\n", err.Error())
			return
		}
		letter, exists := letterMap[letterRune]
		if !exists {
			newLetter := &Letter{
				character: letterRune,
				rowCount:  globalRowCount,
				rowWidth:  rowWidth,
				rows:      make([]string, globalRowCount),
			}
			newLetter.rows[rowNum-1] = rowString
			letterMap[letterRune] = newLetter
		} else {
			if rowWidth != letter.rowWidth {
				fmt.Printf("Error: Length of string value inconsistant within letter ('%s', %d vs %d)\n", string(letterRune), rowWidth, letter.rowWidth)
				return
			}
			if letter.rows[rowNum-1] != "" {
				fmt.Printf("Error: Duplicate declaration of row ('%s', row %d)\n", string(letterRune), rowNum)
				return
			}
			letter.rows[rowNum-1] = rowString
		}
	}

	if *lettersPerRow == 0 || *lettersPerRow < -1 {
		fmt.Printf("Error: letters per row must be positive integer or -1\n")
		return
	}

	fmt.Printf("Input: %s\n", rawInput)
	input := strings.ToLower(rawInput)

	totalLetters := []*Letter{}
	for _, letterRune := range input {
		letter, exists := letterMap[letterRune]
		if !exists {
			fmt.Printf("No entry for rune '%s' in map\n", string(letterRune))
			return
		}
		totalLetters = append(totalLetters, letter)
	}

	if *lettersPerRow == -1 {
		*lettersPerRow = len(totalLetters)
	}

	trueLetters := [][]*Letter{}
	trueLetters = append(trueLetters, []*Letter{})
	lineLength := 0
	for _, letter := range totalLetters {
		if string(letter.character) == *splitCharacter {
			trueLetters = append(trueLetters, []*Letter{})
			fmt.Printf("Saw %s, appending new slice\n", string(letter.character))
			fmt.Printf("Old size: %d New Size: %d\n", len(trueLetters)-1, len(trueLetters))
			lineLength = 0
		} else if lineLength == *lettersPerRow {
			trueLetters = append(trueLetters, []*Letter{})
			fmt.Printf("Saw %s, appending new slice\n", string(letter.character))
			fmt.Printf("Old size: %d New Size: %d\n", len(trueLetters)-1, len(trueLetters))
			lineLength = 0
			trueLetters[len(trueLetters)-1] = append(trueLetters[len(trueLetters)-1], letter)
			fmt.Printf("Appending %s to trueLetters at %d\n", string(letter.character), len(trueLetters)-1)
			lineLength++
		} else {
			trueLetters[len(trueLetters)-1] = append(trueLetters[len(trueLetters)-1], letter)
			fmt.Printf("Appending %s to trueLetters at %d\n", string(letter.character), len(trueLetters)-1)
			lineLength++
		}
	}

	for i := 0; i < len(trueLetters); i++ {
		for rowNum := 0; rowNum < globalRowCount; rowNum++ {
			for _, row := range trueLetters[i] {
				fmt.Print(row.rows[rowNum])
			}
			fmt.Println()
		}
	}

}
