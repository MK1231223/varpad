package variablelanguage

import (
	"errors"
	"slices"
	"strconv"
	"varpad/internal/constants"
)

//WARNING!
//I HATE THIS PARSER SO MUCH
//EVERY SINGLE LINE OF THIS CODE SHOULD BE THROWN IN TO GARBAGE
//BUT I DECIDED TO NOT TO CHANGE ANY OF THE FUCKING LINE IN THIS CODE HERE BECAUSE OF THESE BITCHASS DEADLINES
//MY HEAD IS BOILING LIKE A GREASE ON A HOT PAN
//well... i will rewrite it.... someday.... BUT I'M NOT REWRITING THIS PIECE OF WORKING GARBAGE TODAY!!! FUCK YOU

//17.04.2026

func Parse(Tokens *[]Token, StringVariableBuffer *[]StringValue, IntegerVariableBuffer *[]IntegerValue) error {
	equasion := false
	arguments := []Token{}
	variable := Token{}
	line := 0

	searchThroughStringVariables := func(token *Token) (int, int, string, error) {
		for idx, element := range *StringVariableBuffer {
			if element.Name == token.Value {
				return element.Line, idx, element.Val, nil
			}
		}
		return -1, -1, "", errors.New("var not found")
	}

	searchThroughIntegerVariables := func(token *Token) (int, int, int, error) {
		for idx, element := range *IntegerVariableBuffer {
			if element.Name == token.Value {
				return element.Line, idx, element.Val, nil
			}
		}
		return -1, -1, 0, errors.New("var not found")
	}

	invalidSynthaxErrorCleanUp := func() {
		if variable.Value != "" {
			if _, idxStr, _, err := searchThroughStringVariables(&variable); err == nil {
				*StringVariableBuffer = slices.Delete(*StringVariableBuffer, idxStr, idxStr+1)
				return
			}
			if _, idxStr, _, err := searchThroughIntegerVariables(&variable); err == nil {
				*IntegerVariableBuffer = slices.Delete(*IntegerVariableBuffer, idxStr, idxStr+1)
				return
			}
		}
	}
	

	parseNewVariables := func (found *bool) error {
		l, idxInt, _, err := searchThroughIntegerVariables(&variable)
		if err == nil {
			if l == line {
				if len(arguments) > 0 {
					str, _ := strconv.Atoi(arguments[0].Value)
					(*IntegerVariableBuffer)[idxInt].Val = str
				} else {
					*IntegerVariableBuffer = slices.Delete(*IntegerVariableBuffer, idxInt, idxInt+1)
					return errors.New("no value!")
				}
				*found = true
			} else {
				return errors.New("variable already exists!")
			}
		}
		l, idxStr, _, err := searchThroughStringVariables(&variable)
		if err == nil {
			if l == line {
				if len(arguments) > 0 {
					(*StringVariableBuffer)[idxStr].Val = arguments[0].Value
				} else {
					*StringVariableBuffer = slices.Delete(*StringVariableBuffer, idxStr, idxStr+1)
					return errors.New("no value!")
				}
				*found = true
			} else {
				return errors.New("variable already exists!")
			}
		}
		if !*found {
			if len(arguments) > 0 {
				if arguments[0].Token_type == constants.String {
					*StringVariableBuffer = append(*StringVariableBuffer, StringValue{Name: variable.Value, Val: arguments[0].Value, Line: line})
				}
				if arguments[0].Token_type == constants.Integer {
					res, _ := strconv.Atoi(arguments[0].Value)
					*IntegerVariableBuffer = append(*IntegerVariableBuffer, IntegerValue{Name: variable.Value, Val: res, Line: line})
				}
			} else {
				return errors.New("no value!")
			}
		}
		return nil
	}

	parseExistingVariables := func (found *bool) error {
		if _, idxInt, _, err := searchThroughIntegerVariables(&variable); err == nil {
			if len(arguments) > 0 {
				if arguments[0].Token_type != constants.Integer {
					return errors.New("can't change to value of different type!")
				}
				str, _ := strconv.Atoi(arguments[0].Value)
				(*IntegerVariableBuffer)[idxInt].Val = str
			} else {
				return errors.New("no value!")
			}
			*found = true
		}
		if _, idxStr, _, err := searchThroughStringVariables(&variable); err == nil {
			if len(arguments) > 0 {
				if arguments[0].Token_type != constants.String {
					return errors.New("can't change to value of different type!")
				}
				(*StringVariableBuffer)[idxStr].Val = arguments[0].Value
			} else {
				return errors.New("no value!")
			}
			*found = true
		}
		if !*found {
			return errors.New("variable does not exist!")
		}
		return nil
	}

	parseArgs := func() error {
		//SEARCH VARS
		for idx, arg := range arguments {
			switch arg.Token_type {
			case constants.NewVariable:
				return errors.New("")
			case constants.ExistingVariable:
				_, _, strResult, err := searchThroughStringVariables(&arg)
				if err == nil {
					arguments[idx] = Token{Token_type: constants.String, Value: strResult}
					continue
				} else {
					invalidSynthaxErrorCleanUp()
				}
				_, _, intResult, err := searchThroughIntegerVariables(&arg)
				if err == nil {
					arguments[idx] = Token{Token_type: constants.Integer, Value: strconv.Itoa(intResult)}
					continue
				} else {
					invalidSynthaxErrorCleanUp()
				}
				return errors.New(err.Error())
			}
		}
		//SEARCH MULTIPLICATION AND DIVISION
		for i := 0; i < len(arguments); {
			arg := arguments[i]
			if arg.Token_type == constants.Multiply || arg.Token_type == constants.Divide {
				if i == 0 || i == len(arguments)-1 {
					invalidSynthaxErrorCleanUp()
					return errors.New("invalid syntax")
				}
				op1 := arguments[i-1]
				op2 := arguments[i+1]
				if op1.Token_type == constants.String || op2.Token_type == constants.String {
					invalidSynthaxErrorCleanUp()
					return errors.New("invalid syntax")
				}
				num1, _ := strconv.Atoi(op1.Value)
				num2, _ := strconv.Atoi(op2.Value)
				var result int
				if arg.Token_type == constants.Multiply {
					result = num1 * num2
				} else {
					result = num1 / num2
				}
				arguments = slices.Replace(arguments, i-1, i+2, Token{Token_type: constants.Integer, Value: strconv.Itoa(result)})
				i--
				if i < 0 {
					i = 0
				}
				continue
			}
			i++
		}
		//SEARCH SUBTRACTION AND ADDITION
		for i := 0; i < len(arguments); {
			arg := arguments[i]
			if arg.Token_type == constants.Plus || arg.Token_type == constants.Minus {
				if i == 0 || i == len(arguments)-1 {
					invalidSynthaxErrorCleanUp()
					return errors.New("invalid syntax")
				}
				op1 := arguments[i-1]
				op2 := arguments[i+1]
				// PLUS
				if arg.Token_type == constants.Plus {
					// string concatenation if any operand is string
					if op1.Token_type == constants.String || op2.Token_type == constants.String {
						arguments = slices.Replace(arguments, i-1, i+2, Token{Token_type: constants.String, Value: op1.Value + op2.Value})
						if _, idx, _, err := searchThroughIntegerVariables(&variable); err == nil {
							*IntegerVariableBuffer = slices.Delete(*IntegerVariableBuffer, idx, idx+1)
						}
					} else {
						num1, _ := strconv.Atoi(op1.Value)
						num2, _ := strconv.Atoi(op2.Value)
						arguments = slices.Replace(arguments, i-1, i+2, Token{Token_type: constants.Integer, Value: strconv.Itoa(num1 + num2)})
					}
				}
				// MINUS
				if arg.Token_type == constants.Minus {
					if op1.Token_type == constants.String || op2.Token_type == constants.String {
						invalidSynthaxErrorCleanUp()
						return errors.New("invalid syntax")
					}
					num1, _ := strconv.Atoi(op1.Value)
					num2, _ := strconv.Atoi(op2.Value)
					arguments = slices.Replace(arguments, i-1, i+2, Token{Token_type: constants.Integer, Value: strconv.Itoa(num1 - num2)})
				}
				i--
				if i < 0 {
					i = 0
				}
				continue
			}
			i++
		}

		found := false

		if variable.Token_type == constants.ExistingVariable {
			err := parseExistingVariables(&found)
			if err != nil {
				return errors.New(err.Error())
			}
		} else {
			err := parseNewVariables(&found)
			if err != nil {
				return errors.New(err.Error())
			}
		}

		equasion = false
		arguments = []Token{}
		variable = Token{}
		line++
		return nil
	}

	for _, token := range *Tokens {
		addToArguments := func() error {
			if !equasion {
				return errors.New("")
			}
			arguments = append(arguments, token)
			return nil
		}

		writeVariable := func() error {
			if (variable == Token{}) {
				variable = token
				equasion = true
			} else {
				err := addToArguments()
				if err != nil {
					return errors.New(err.Error())
				}
			}
			return nil
		}

		switch token.Token_type {
		case constants.ExistingVariable, constants.NewVariable:
			err := writeVariable()
			if err != nil {
				return errors.New(err.Error())
			}
		case constants.Equals:
			if equasion && len(arguments) > 0 {
				return errors.New("Invalid synthax!")
			}
			equasion = true
		case constants.Newline:
			if variable.Value == "" {
				return nil
			}
			err := parseArgs()
			if err != nil {
				return errors.New(err.Error())
			}
		case constants.Garbage:
			invalidSynthaxErrorCleanUp()
			return errors.New("Invalid synthax!")
		default:
			err := addToArguments()
			if err != nil {
				return errors.New(err.Error())
			}
		}
	}

	return nil
}
