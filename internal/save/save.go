package save

import (
	"errors"
	"os"
	"strconv"
	"strings"
	app_math "varpad/internal/math"
	variablelanguage "varpad/internal/variable_language"
)

func Save(TextBuffer *[][]rune, StringVariableBuffer *[]variablelanguage.StringValue, IntegerVariableBuffer *[]variablelanguage.IntegerValue, VariableBlockRange app_math.Vector1) {
	filename := "export"
	fileformat := ".txt"
	file, err := os.OpenFile("export.txt", os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
	if err != nil {
		for count := 1; ; count++ {
			file, err = os.OpenFile(filename+strconv.Itoa(count)+fileformat, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
			if err == nil {
				break
			}
		}
	}

	searchThroughStringVariables := func(word string) (string, error) {
		for _, element := range *StringVariableBuffer {
			if element.Name == word {
				return element.Val, nil
			}
		}
		return "", errors.New("var not found")
	}

	searchThroughIntegerVariables := func(word string) (int, error) {
		for _, element := range *IntegerVariableBuffer {
			if element.Name == word {
				return element.Val, nil
			}
		}
		return 0, errors.New("var not found")
	}

	line := ""
	lineSplit := []string{}
	for lIdx, l := range *TextBuffer {
		if lIdx <= VariableBlockRange.X2 {
			continue
		}

		line = string(l)
		lineSplit = strings.Split(line, " ")

		for idx, word := range lineSplit {
			if len(word) == 0 {
				continue
			}
			if word[0] == '%' {
				if val, err := searchThroughStringVariables(word[1:]); err == nil {
					lineSplit[idx] = val
					break
				}
				if val, err := searchThroughIntegerVariables(word[1:]); err == nil {
					lineSplit[idx] = strconv.Itoa(val)
					break
				}
			}
		}

		if lIdx != len(*TextBuffer)-1 {
			lineSplit = append(lineSplit, "\n")
		}

		file.WriteString(strings.Join(lineSplit, " "))
	}
	file.Close()
}
