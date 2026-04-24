package elements

import (
	"strconv"

	constants "varpad/internal/constants"
	app_math "varpad/internal/math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Wndw struct {
	W     int32
	H     int32
	MinW  int32
	MinH  int32
	Title string
	FPS   int32
}

func (wndw *Wndw) Update() {
	wndw.W = int32(rl.GetRenderWidth())
	wndw.H = int32(rl.GetRenderHeight())
}

func (wndw *Wndw) UpdateTitle(s string) {
	wndw.Title = s
	rl.SetWindowTitle(wndw.Title)
}

type Line struct {
	Text string
	X    int32
	Y    int32
}

type WritingSpaceSettings struct {
	//TEMP
	SelectedLine         *int
	SelectedCharacterNum *int
	TextBuffer           *[][]rune
	VariableBlockRange   app_math.Vector1

	HasVariableBlock *bool
	CodeHasAnError   *bool

	DrawFromLine      *int
	DrawFromCharacter *int

	ScrollYLocked *bool
	ScrollXLocked *bool

	MousePos *rl.Vector2

	Wrtsp *rl.RectangleInt32

	//CONST
	TextSize int32

	EvenLineColor     rl.Color
	OddLineColor      rl.Color
	TextColor         rl.Color
	VarTextColor      rl.Color
	BorderColor       rl.Color
	VarErrorTextColor rl.Color

	TextFont *rl.Font
}

// WARNING! SHITASS CODE!
func LayoutSpaceLines(Settings *WritingSpaceSettings) {
	var (
		settings = *Settings

		wrtsp         = settings.Wrtsp
		scrollYLocked = settings.ScrollYLocked
		font          = settings.TextFont
		mousePos      = settings.MousePos

		lineNumber             int = *settings.DrawFromLine
		lineNumberMarginOffset int = 0
		pixels                 int = 0
		lineColor              rl.Color
		textColor              rl.Color

		textBackShift int //USED FOR HORIZONTAL TEXT SCROLLING
	)

	lastNumberLengh := len(strconv.Itoa(len(*settings.TextBuffer)))
	lineNumberMarginOffset = lastNumberLengh * int(float32(settings.TextSize)*0.45)

	textBackShift = (*settings.DrawFromCharacter) * int(float32(settings.TextSize)*0.45)

	drawLine := func() {
		line := Line{X: wrtsp.X - 1, Y: wrtsp.Y + int32(pixels)}
		isLineNumberEven := lineNumber%2 == 0
		if lineNumber < len(*settings.TextBuffer) {
			line.Text = string((*settings.TextBuffer)[lineNumber])
		}

		func() { //INIT LINES COLOR
			if isLineNumberEven {
				lineColor = settings.EvenLineColor
			} else {
				lineColor = settings.OddLineColor
			}
		}()

		func() { //INIT TEXT COLOR
			if lineNumber >= settings.VariableBlockRange.X1 && lineNumber <= settings.VariableBlockRange.X2 {
				if *settings.CodeHasAnError {
					textColor = settings.VarErrorTextColor
					return
				}
				textColor = settings.VarTextColor
			} else {
				textColor = settings.TextColor
			}
		}()

		//LINE RECTANGLE
		lineRectangle := rl.RectangleInt32{
			X:      line.X,
			Y:      line.Y,
			Width:  wrtsp.Width + 1,
			Height: settings.TextSize,
		}
		rl.DrawRectangle(
			lineRectangle.X,
			lineRectangle.Y,
			lineRectangle.Width,
			lineRectangle.Height,
			lineColor,
		)

		//LINE NUMBER
		rl.DrawTextEx(
			*font,
			strconv.Itoa(lineNumber+1),
			rl.Vector2{
				X: float32(line.X + constants.TextMargin),
				Y: float32(wrtsp.Y + int32(pixels)),
			},
			float32(settings.TextSize),
			0,
			settings.TextColor,
		)

		func() { //LINE TEXT
			textLengh := float32(len((*settings.TextBuffer)[lineNumber])) * float32(settings.TextSize) * 0.45
			textX := float32(wrtsp.X + int32(lineNumberMarginOffset) + constants.TextMargin*3)
			textY := float32(wrtsp.Y + int32(pixels))
			rl.BeginScissorMode(int32(textX), int32(textY), int32(textLengh), settings.TextSize)
			rl.DrawTextEx(
				*font,
				line.Text,
				rl.Vector2{
					X: textX - float32(textBackShift),
					Y: textY,
				},
				float32(settings.TextSize),
				0,
				textColor,
			)
			rl.EndScissorMode()
			lineNumber++
			pixels += int(settings.TextSize)
		}()

		func() { //DRAW SELECTION LINE
			if *settings.SelectedLine == lineNumber {
				selectionLineX := float32(*settings.SelectedCharacterNum-*settings.DrawFromCharacter)*(float32(settings.TextSize)*0.45) + float32(lineNumberMarginOffset+constants.TextMargin*5)
				startPos := rl.Vector2{
					X: selectionLineX,
					Y: float32(line.Y) + float32(app_math.PercentageI32(settings.TextSize, 5)),
				}
				endPos := rl.Vector2{
					X: selectionLineX,
					Y: float32(line.Y) + float32(app_math.PercentageI32(settings.TextSize, 95)),
				}
				rl.DrawLineEx(startPos, endPos, 3, settings.BorderColor)
			}
		}()

		//MOUSE INPUT
		isMouseHovering := rl.CheckCollisionPointRec(*mousePos, lineRectangle.ToFloat32())
		if isMouseHovering {
			if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
				*settings.SelectedLine = lineNumber
				if !*settings.ScrollXLocked {
					*settings.SelectedCharacterNum = 0
					return
				}
				*settings.SelectedCharacterNum = len(line.Text)
			}
		}
	}

	drawSeparatingLine := func() {
		separatingLineX := int(wrtsp.X) + lineNumberMarginOffset + constants.TextMargin*2
		rl.DrawLine(int32(separatingLineX), wrtsp.Y, int32(separatingLineX), wrtsp.Height+settings.TextSize, settings.BorderColor)
	}

	checkOverflowingText := func() {
		matches := 0
		for line := 0; line < len(*settings.TextBuffer); line++ {
			textLine := (*settings.TextBuffer)[line]
			textLengh := (float32(len(textLine)) * float32(settings.TextSize) * 0.45) - float32(textBackShift)
			if textLengh+constants.TextMargin*5 > float32(wrtsp.Width) {
				matches++
				*settings.ScrollXLocked = false
			}
		}
		if matches == 0 {
			*settings.ScrollXLocked = true
		}
	}

	for {
		checkOverflowingText()
		if (pixels + int(settings.TextSize)) > int(wrtsp.Height) {
			drawLine()
			*scrollYLocked = false
			drawSeparatingLine()
			break
		}
		*scrollYLocked = true
		drawLine()
		drawSeparatingLine()
		if lineNumber >= len(*settings.TextBuffer) {
			break
		}
	}
}
