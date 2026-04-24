package variablelanguage

import (
	"unicode"
	"varpad/internal/constants"
	app_math "varpad/internal/math"
)

func FileHasVariableBlock(TextBuffer *[][]rune) (bool, app_math.Vector1) {
	hasVariables := false
	variableBlockPosition := app_math.Vector1{X1: -1, X2: -1}

	for number, line := range *TextBuffer {
		if number == 0 {
			if string(line) == constants.VarStart {
				variableBlockPosition.X1 = 0
			} else {
				variableBlockPosition.ResetTo(-1)
				return hasVariables, variableBlockPosition
			}
		}

		if string(line) == constants.VarEnd {
			variableBlockPosition.X2 = number
		}
	}

	if variableBlockPosition.X2 == -1 {
		variableBlockPosition.ResetTo(-1)
		return hasVariables, variableBlockPosition
	} else {
		hasVariables = true
	}

	return hasVariables, variableBlockPosition
}

func Split(CodeBuffer *[][]rune) []string {
	result := []string{}
	word := ""
	addToResult := func() {
		if len(word) == 0 {
			return
		}
		result = append(result, word)
		word = ""
	}
	quotes := 0

	for _, line := range *CodeBuffer {
		emptySymbols := 0
		for idx, character := range line {
			if (character == ' ' || character == '\t') && ((quotes%2 == 0) || (quotes == 0)) {
				emptySymbols++
				if idx != 0 {
					addToResult()
				}
				if idx == len(line)-1 && emptySymbols != len(line) {
					result = append(result, "\n")
				}
			} else {
				word += string(character)
				if idx == len(line)-1 {
					addToResult()
					result = append(result, "\n")
				}
			}
			if character == '"' {
				quotes++
			}
		}
	}
	return result
}

func Tokenize(SplitBuffer *[]string) []Token {
	isNewVariable := func(s string) bool {
		res := false
		if len(s) == 1 {
			return false
		}
		identifierExists := false
		for idx, character := range s {
			if character == '$' && idx == 0 {
				identifierExists = true
				continue
			}
			if !((character >= 'a' && character <= 'z') || (character >= 'A' && character <= 'Z')) {
				return false
			}
		}
		if identifierExists {
			res = true
		}
		return res
	}

	isExistingVariable := func(s string) bool {
		res := false
		if len(s) == 1 {
			return false
		}
		identifierExists := false
		for idx, character := range s {
			if character == '%' && idx == 0 {
				identifierExists = true
				continue
			}
			if !((character >= 'a' && character <= 'z') || (character >= 'A' && character <= 'Z')) {
				return false
			}
		}
		if identifierExists {
			res = true
		}
		return res
	}

	isString := func(s string) bool {
		res := false
		if len(s) < 2 {
			return false
		}

		if s[0] == '"' && s[len(s)-1] == '"' {
			res = true
		}

		for idx, character := range s {
			if character == '"' && idx != 0 && idx != len(s)-1 {
				return false
			}
		}

		return res
	}

	isInteger := func(s string) bool {
		res := false
		if len(s) == 1 {
			if unicode.IsPunct(rune(s[0])) || unicode.IsSymbol(rune(s[0])) || unicode.IsLetter(rune(s[0])) {
				return false
			}
		}
		numbers := 0
		for _, character := range s {
			if unicode.IsNumber(character) {
				numbers++
			}
		}
		if numbers == len(s) {
			res = true
		}
		return res
	}

	isOperator := func(s string) bool {
		res := false
		if len(s) > 1 {
			return false
		}
		if unicode.IsPunct(rune(s[0])) || unicode.IsSymbol(rune(s[0])) {
			res = true
		}
		return res
	}

	isNewline := func(s string) bool {
		if s == "\n" {
			return true
		}
		return false
	}

	tokens := []Token{}

	for _, word := range *SplitBuffer {
		if isNewline(word) {
			tokens = append(tokens, Token{Token_type: constants.Newline})
			continue
		}
		if isNewVariable(word) {
			tokens = append(tokens, Token{Token_type: constants.NewVariable, Value: word[1:]})
			continue
		}
		if isExistingVariable(word) {
			tokens = append(tokens, Token{Token_type: constants.ExistingVariable, Value: word[1:]})
			continue
		}
		if isOperator(word) {
			switch word {
			case "+":
				tokens = append(tokens, Token{Token_type: constants.Plus})
			case "-":
				tokens = append(tokens, Token{Token_type: constants.Minus})
			case "*":
				tokens = append(tokens, Token{Token_type: constants.Multiply})
			case "/":
				tokens = append(tokens, Token{Token_type: constants.Divide})
			case "=":
				tokens = append(tokens, Token{Token_type: constants.Equals})
			default:
				tokens = append(tokens, Token{Token_type: constants.Garbage, Value: word})
			}
			continue
		}
		if isInteger(word) {
			tokens = append(tokens, Token{Token_type: constants.Integer, Value: word})
			continue
		}
		if isString(word) {
			tokens = append(tokens, Token{Token_type: constants.String, Value: word[1 : len(word)-1]})
			continue
		}
		tokens = append(tokens, Token{Token_type: constants.Garbage, Value: word})
	}

	return tokens
}
