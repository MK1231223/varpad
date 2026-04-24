package app_text

import (
	app_math "varpad/internal/math"
	variablelanguage "varpad/internal/variable_language"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type ProcessFileTextSettings struct {
	TextBuffer            *[][]rune
	CodeBuffer            *[][]rune
	HasVariableBlock      *bool
	VariableBlockRange    *app_math.Vector1
	CodeHasAnError        *bool
	StringVariableBuffer  *[]variablelanguage.StringValue
	IntegerVariableBuffer *[]variablelanguage.IntegerValue
	UpperPannel           *rl.RectangleInt32
	ErrorColor            rl.Color
	Font                  *rl.Font
}

func ProcessFileText(settings *ProcessFileTextSettings) {
	*settings.HasVariableBlock, *settings.VariableBlockRange = variablelanguage.FileHasVariableBlock(settings.TextBuffer)
	if *settings.HasVariableBlock {
		updateCodeBuffer(settings)
		splitbuffer := variablelanguage.Split(settings.CodeBuffer)
		tokenBuffer := variablelanguage.Tokenize(&splitbuffer)
		err := variablelanguage.Parse(&tokenBuffer, settings.StringVariableBuffer, settings.IntegerVariableBuffer)
		if err != nil {
			*settings.CodeHasAnError = true
			rl.DrawTextEx(*settings.Font, err.Error(), rl.Vector2{X: float32(settings.UpperPannel.X), Y: float32(settings.UpperPannel.Y)}, 25, 0, settings.ErrorColor)
		} else {
			*settings.CodeHasAnError = false
		}
	}
}

func updateCodeBuffer(settings *ProcessFileTextSettings) {
	*settings.CodeBuffer = [][]rune{}
	for _, line := range (*settings.TextBuffer)[settings.VariableBlockRange.X1+1 : settings.VariableBlockRange.X2] {
		*settings.CodeBuffer = append(*settings.CodeBuffer, line)
	}
}
