package variablelanguage

type Token struct {
	Token_type string
	Value      string
}

type IntegerValue struct {
	Name string
	Val  int
	Line int
}

type StringValue struct {
	Name string
	Val  string
	Line int
}
