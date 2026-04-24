package constants

const (
	//VALUES
	DYNAMIC = -1
	NULL    = 0

	//LAYOUT
	WritingSpaceMargin = 10
	TextMargin         = 5

	//FILE
	FileFormat = ".vrpd"

	//STRING
	MainTitle          = "VARPAD"
	NewProjectTitle    = MainTitle + " - NEW PROJECT* " + NativeFormatOpened
	TxtOpened          = "(.txt)"
	NativeFormatOpened = "(" + FileFormat + ")"
)

const (
	//LANGUAGE
	VarStart = "||VAR:"
	VarEnd   = ":END||"

	//TOKEN TYPES
	Integer          = "INTEGER"
	String           = "STRING"
	NewVariable      = "NEW-VARIABLE"
	ExistingVariable = "EXISTING-VARIABLE"
	Newline          = "NEWLINE"
	Garbage          = "GARBAGE"
	Plus             = "PLUS"
	Minus            = "MINUS"
	Multiply         = "MULTIPLY"
	Divide           = "DIVIDE"
	Equals           = "EQUALS"
)

const (
	//PROJECT OPEN FORMAT
	TxtFormat = iota
	NativeFormat
)
