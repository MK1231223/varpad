package input

import (
	app_math "varpad/internal/math"
	"varpad/internal/save"
	variablelanguage "varpad/internal/variable_language"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type InputHandler struct{}
type InputHandlerSettings struct {
	SelectedLine          *int
	TextBuffer            *[][]rune
	SelectedCharacterNum  *int
	PointPosition         rl.Vector2
	WritingSpaceOutline   *rl.RectangleInt32
	ScrollPower           int
	ScrollY               *int
	ScrollYLocked         *bool
	ScrollX               *int
	ScrollXLocked         *bool
	StringVariableBuffer  *[]variablelanguage.StringValue
	IntegerVariableBuffer *[]variablelanguage.IntegerValue
	VariableBlockRange    app_math.Vector1
}

// TEXT
// -------------------------------------------------------------------------------------------------------------
func (ih InputHandler) handleEnter(settings *InputHandlerSettings) {
	if *settings.SelectedLine == 0 {
		return
	}
	if !rl.IsKeyPressed(rl.KeyEnter) && !rl.IsKeyPressedRepeat(rl.KeyEnter) {
		return
	}
	currentLine := (*settings.TextBuffer)[*settings.SelectedLine-1]
	lineBuffer := []rune{}
	if *settings.SelectedCharacterNum < len(currentLine) {
		lineBuffer = currentLine[*settings.SelectedCharacterNum:]
		(*settings.TextBuffer)[*settings.SelectedLine-1] = currentLine[:*settings.SelectedCharacterNum]
	}
	*settings.SelectedCharacterNum = 0
	if len(*settings.TextBuffer) > *settings.SelectedLine {
		*settings.TextBuffer = append((*settings.TextBuffer)[:*settings.SelectedLine], append([][]rune{lineBuffer}, (*settings.TextBuffer)[*settings.SelectedLine:]...)...)
		*settings.SelectedLine++
		return
	}
	*settings.TextBuffer = append(*settings.TextBuffer, lineBuffer)
	*settings.SelectedLine++
}
func (ih InputHandler) handleBackspace(settings *InputHandlerSettings) {
	if *settings.SelectedLine == 0 {
		return
	}
	if !rl.IsKeyPressed(rl.KeyBackspace) && !rl.IsKeyPressedRepeat(rl.KeyBackspace) {
		return
	}

	bufferPosition := *settings.SelectedLine - 1

	rawLine := &(*settings.TextBuffer)[bufferPosition]

	if len(*rawLine) > 0 && *settings.SelectedCharacterNum > 0 {
		if len(*rawLine) > *settings.SelectedCharacterNum {
			*rawLine = append((*rawLine)[:*settings.SelectedCharacterNum-1], (*rawLine)[*settings.SelectedCharacterNum:]...)
		} else {
			*rawLine = (*rawLine)[:len(*rawLine)-1]
		}
		*settings.SelectedCharacterNum--
		return
	}

	setCursor := func() {
		*settings.SelectedCharacterNum = len((*settings.TextBuffer)[bufferPosition-1])
	}

	isMovingTextToPrevLine := len(*rawLine) > 0 && *settings.SelectedCharacterNum == 0 && *settings.SelectedLine > 1

	if isMovingTextToPrevLine {
		setCursor()
		(*settings.TextBuffer)[bufferPosition-1] = append((*settings.TextBuffer)[bufferPosition-1], *rawLine...)
		isMovingTextToPrevLine = true
	}

	if *settings.SelectedLine <= 1 {
		return
	}

	*settings.SelectedLine--
	*settings.TextBuffer = append((*settings.TextBuffer)[:*settings.SelectedLine], (*settings.TextBuffer)[*settings.SelectedLine+1:]...)
	if !isMovingTextToPrevLine {
		setCursor()
	}
}
func (ih InputHandler) handleLetters(settings *InputHandlerSettings) {
	if *settings.SelectedLine == 0 {
		return
	}
	key := rl.GetCharPressed()
	addCharacter := func() {
		rawLine := &(*settings.TextBuffer)[*settings.SelectedLine-1]

		if *settings.SelectedCharacterNum < len(*rawLine) {
			*rawLine = append((*rawLine)[:*settings.SelectedCharacterNum], append([]rune{rune(key)}, (*rawLine)[*settings.SelectedCharacterNum:]...)...)
		} else {
			*rawLine = append(*rawLine, rune(key))
		}

		*settings.SelectedCharacterNum++
	}
	for key > 0 {
		r := rune(key)

		if r != '\n' && r != '\t' && r != '\r' {
			addCharacter()
		}

		key = rl.GetCharPressed()
	}
}
func (ih InputHandler) handleTAB(settings *InputHandlerSettings) {
	if !rl.IsKeyPressed(rl.KeyTab) && !rl.IsKeyPressedRepeat(rl.KeyTab) {
		return
	}
	if *settings.SelectedLine == 0 {
		return
	}
	bufferPosition := *settings.SelectedLine - 1
	rawLine := &(*settings.TextBuffer)[bufferPosition]

	if len(*rawLine) > *settings.SelectedCharacterNum {
		*rawLine = append((*rawLine)[:*settings.SelectedCharacterNum], append([]rune{'\t'}, (*rawLine)[*settings.SelectedCharacterNum:]...)...)
	} else {
		*rawLine = append(*rawLine, '\t')
	}
	*settings.SelectedCharacterNum++
}
func (ih InputHandler) handleClipboard(settings *InputHandlerSettings) {
	if *settings.SelectedLine == 0 {
		return
	}
	paste := func() {
		currentLine := (*settings.TextBuffer)[*settings.SelectedLine-1]
		copiedText := rl.GetClipboardText()
		if copiedText == "" {
			return
		}
		convertedCopiedText := []rune(copiedText)
		if *settings.SelectedCharacterNum < len(currentLine) {
			bufferPosition := *settings.SelectedLine - 1
			line := &((*settings.TextBuffer)[bufferPosition])
			*line = append(
				(*line)[:*settings.SelectedCharacterNum],
				append(convertedCopiedText, (*line)[*settings.SelectedCharacterNum:]...)...,
			)
			*settings.SelectedCharacterNum += len(convertedCopiedText)
			return
		}
		*settings.SelectedCharacterNum += len(convertedCopiedText)
		(*settings.TextBuffer)[*settings.SelectedLine-1] = append((*settings.TextBuffer)[*settings.SelectedLine-1], convertedCopiedText...)
	}
	if rl.IsKeyDown(rl.KeyLeftControl) || rl.IsKeyDown(rl.KeyRightControl) {
		if rl.IsKeyPressed(rl.KeyV) {
			paste()
		}
	}
}
func (ih InputHandler) TextInput(settings *InputHandlerSettings) {
	ih.handleEnter(settings)
	ih.handleBackspace(settings)
	ih.handleLetters(settings)
	ih.handleTAB(settings)
	ih.handleClipboard(settings)
}

// ARROWS
// -------------------------------------------------------------------------------------------------------------
func (ih InputHandler) ArrowInput(settings *InputHandlerSettings) {
	bufferPosition := *settings.SelectedLine - 1
	adjustCharacterSelection := func(key int) {
		switch key {
		case rl.KeyUp:
			if *settings.SelectedCharacterNum < len((*settings.TextBuffer)[bufferPosition-1]) {
				return
			}
			*settings.SelectedCharacterNum = len((*settings.TextBuffer)[bufferPosition-1])
		case rl.KeyDown:
			if *settings.SelectedCharacterNum < len((*settings.TextBuffer)[bufferPosition+1]) {
				return
			}
			*settings.SelectedCharacterNum = len((*settings.TextBuffer)[bufferPosition+1])
		}
	}

	if (rl.IsKeyPressed(rl.KeyUp) || rl.IsKeyPressedRepeat(rl.KeyUp)) && *settings.SelectedLine > 0 {
		if *settings.SelectedLine == 1 {
			return
		}
		adjustCharacterSelection(rl.KeyUp)
		*settings.SelectedLine--
	}
	if (rl.IsKeyPressed(rl.KeyDown) || rl.IsKeyPressedRepeat(rl.KeyDown)) && *settings.SelectedLine > 0 {
		if len(*settings.TextBuffer) == *settings.SelectedLine {
			return
		}
		adjustCharacterSelection(rl.KeyDown)
		*settings.SelectedLine++
	}
	if (rl.IsKeyPressed(rl.KeyRight) || rl.IsKeyPressedRepeat(rl.KeyRight)) && *settings.SelectedLine > 0 {
		if len((*settings.TextBuffer)[bufferPosition]) == *settings.SelectedCharacterNum {
			if len(*settings.TextBuffer) == *settings.SelectedLine {
				return
			}
			*settings.SelectedLine++
			*settings.SelectedCharacterNum = 0
			return
		}
		*settings.SelectedCharacterNum++
	}
	if (rl.IsKeyPressed(rl.KeyLeft) || rl.IsKeyPressedRepeat(rl.KeyLeft)) && *settings.SelectedLine > 0 {
		if *settings.SelectedCharacterNum == 0 {
			if *settings.SelectedLine == 1 {
				return
			}
			*settings.SelectedLine--
			*settings.SelectedCharacterNum = len((*settings.TextBuffer)[bufferPosition-1])
			return
		}

		*settings.SelectedCharacterNum--
	}

	if *settings.SelectedLine > len(*settings.TextBuffer) {
		*settings.SelectedLine = len(*settings.TextBuffer)
	}
}

// SYSTEM HOTKEYS
// -------------------------------------------------------------------------------------------------------------
func (ih InputHandler) SystemInput(settings *InputHandlerSettings) {
	if rl.IsKeyPressed(rl.KeyF10) {
		save.Save(settings.TextBuffer, settings.StringVariableBuffer, settings.IntegerVariableBuffer, settings.VariableBlockRange)
	}
}

// SCROLL
// -------------------------------------------------------------------------------------------------------------
func (ih InputHandler) ScrollInput(settings *InputHandlerSettings) {
	var ScrollInput float32 = rl.GetMouseWheelMove()
	var ShiftInput bool = rl.IsKeyDown(rl.KeyLeftShift)

	isWritingSpaceHovered := rl.CheckCollisionPointRec(settings.PointPosition, settings.WritingSpaceOutline.ToFloat32())

	if !isWritingSpaceHovered || ScrollInput == 0 {
		return
	}

	if ScrollInput > 0 {
		if ShiftInput {
			*settings.ScrollX -= settings.ScrollPower
		} else {
			*settings.ScrollY -= settings.ScrollPower
		}
	}

	if ScrollInput < 0 {
		if ShiftInput && !*settings.ScrollXLocked {
			*settings.ScrollX += settings.ScrollPower
		} else if !ShiftInput && !*settings.ScrollYLocked {
			*settings.ScrollY += settings.ScrollPower
		}
	}

	if *settings.ScrollY < 0 {
		*settings.ScrollY = 0
	}
	if *settings.ScrollX < 0 {
		*settings.ScrollX = 0
	}
}
