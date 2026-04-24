package main

import (
	constants "varpad/internal/constants"
	elements "varpad/internal/elements"
	"varpad/internal/input"
	app_math "varpad/internal/math"
	app_text "varpad/internal/text"
	variablelanguage "varpad/internal/variable_language"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	projectFormat int

	textSize               int32   = 40
	fontSize               int32   = 96
	scroll_power           int     = 3 //lines
	alllowArrowScrollDelay float32 = 0.5
	arrowScrollDelay       float32 = 0.08

	backgroundColor   = rl.White
	textColor         = rl.Black
	varTextColor      = rl.Blue
	varErrorTextColor = rl.Red
	upbarColor        = rl.Color{R: 210, G: 210, B: 210, A: 255}
	evenLineColor     = rl.Color{R: 230, G: 230, B: 230, A: 255}
	oddLineColor      = rl.White
	borderColor       = rl.Black

	selectedLine         int = 0
	selectedCharacterNum int = 0

	textBuffer         = [][]rune{}
	codeBuffer         = [][]rune{}
	hasVariableBlock   bool
	codeHasAnError     bool
	variableBlockRange app_math.Vector1

	stringVariableBuffer  []variablelanguage.StringValue
	integerVariableBuffer []variablelanguage.IntegerValue
)

func initProject(wndw *elements.Wndw) {
	if len(textBuffer) == 0 {
		textBuffer = [][]rune{[]rune(constants.VarStart), []rune("\t"), []rune(constants.VarEnd)}
		wndw.UpdateTitle(constants.NewProjectTitle)
		projectFormat = constants.NativeFormat
	}
}

func main() {
	//INIT VARS
	var (
		wndw = elements.Wndw{
			W:     800,
			H:     450,
			MinW:  800,
			MinH:  450,
			Title: constants.MainTitle,
			FPS:   240,
		}
		upbarHeight      int32 = 25
		scrollY          int
		scrollX          int
		renderW, renderH int
		pointPosition    rl.Vector2
		scrollYLocked    bool
		scrollXLocked    bool
	)

	//INPUT VARS
	var (
		inputHandler input.InputHandler
	)

	//INIT FONT CODEPOINTS
	codepoints := []int32{}

	for cdpt := int32(32); cdpt < 126; cdpt++ {
		codepoints = append(codepoints, cdpt)
	}

	for cdpt := int32(0x0400); cdpt < 0x0500; cdpt++ {
		codepoints = append(codepoints, cdpt)
	}

	//INIT WINDOW
	rl.SetConfigFlags(rl.FlagWindowResizable)
	rl.InitWindow(wndw.W, wndw.H, wndw.Title)
	rl.SetWindowMinSize(int(wndw.MinW), int(wndw.MinH))
	rl.SetTargetFPS(wndw.FPS)
	rl.SetExitKey(constants.NULL)
	var font = rl.LoadFontEx("etc/fonts/RobotoMono-Regular.ttf", fontSize, codepoints, int32(len(codepoints)))
	rl.SetTextureFilter(font.Texture, rl.FilterBilinear)
	var wrtspSettings = elements.WritingSpaceSettings{}

	initProject(&wndw)

	for !rl.WindowShouldClose() {
		////UPDATE
		//-------------------------------------------------------------------------------------------------------------
		inputHandlerSettings := input.InputHandlerSettings{
			SelectedLine:          &selectedLine,
			TextBuffer:            &textBuffer,
			SelectedCharacterNum:  &selectedCharacterNum,
			PointPosition:         pointPosition,
			WritingSpaceOutline:   nil,
			ScrollPower:           scroll_power,
			ScrollY:               &scrollY,
			ScrollYLocked:         &scrollYLocked,
			ScrollX:               &scrollX,
			ScrollXLocked:         &scrollXLocked,
			StringVariableBuffer:  &stringVariableBuffer,
			IntegerVariableBuffer: &integerVariableBuffer,
			VariableBlockRange:    variableBlockRange,
		}

		upperBar := rl.RectangleInt32{X: 0, Y: 0, Width: int32(renderW), Height: upbarHeight}
		rl.DrawRectangle(upperBar.X, upperBar.Y, upperBar.Width, upbarHeight, upbarColor) //upbar

		processFileTextSettings := app_text.ProcessFileTextSettings{
			TextBuffer:            &textBuffer,
			CodeBuffer:            &codeBuffer,
			HasVariableBlock:      &hasVariableBlock,
			VariableBlockRange:    &variableBlockRange,
			CodeHasAnError:        &codeHasAnError,
			StringVariableBuffer:  &stringVariableBuffer,
			IntegerVariableBuffer: &integerVariableBuffer,
			UpperPannel:           &upperBar,
			ErrorColor:            varErrorTextColor,
			Font:                  &font,
		}

		if projectFormat == constants.NativeFormat {
			app_text.ProcessFileText(&processFileTextSettings)
		}
		inputHandler.SystemInput(&inputHandlerSettings)
		inputHandler.TextInput(&inputHandlerSettings)
		inputHandler.ArrowInput(&inputHandlerSettings)
		renderW = rl.GetRenderWidth()
		renderH = rl.GetRenderHeight()
		pointPosition = rl.GetMousePosition()
		wndw.Update()

		////INIT DRAW
		//-------------------------------------------------------------------------------------------------------------
		rl.BeginDrawing()
		rl.ClearBackground(backgroundColor)

		////MAIN
		//-------------------------------------------------------------------------------------------------------------

		writingSpaceOutline := rl.RectangleInt32{
			X:      constants.WritingSpaceMargin,
			Y:      upbarHeight + constants.WritingSpaceMargin,
			Width:  int32(renderW) - constants.WritingSpaceMargin*2,
			Height: int32(renderH) - (constants.WritingSpaceMargin*2 + upbarHeight),
		}
		wrtspSettings = elements.WritingSpaceSettings{
			TextSize:             textSize,
			EvenLineColor:        evenLineColor,
			OddLineColor:         oddLineColor,
			TextColor:            textColor,
			VarTextColor:         varTextColor,
			BorderColor:          borderColor,
			SelectedLine:         &selectedLine,
			SelectedCharacterNum: &selectedCharacterNum,
			TextBuffer:           &textBuffer,
			VariableBlockRange:   variableBlockRange,
			DrawFromLine:         &scrollY,
			DrawFromCharacter:    &scrollX,
			ScrollYLocked:        &scrollYLocked,
			ScrollXLocked:        &scrollXLocked,
			MousePos:             &pointPosition,
			Wrtsp:                &writingSpaceOutline,
			TextFont:             &font,
			CodeHasAnError:       &codeHasAnError,
			VarErrorTextColor:    varErrorTextColor,
		}
		elements.LayoutSpaceLines(&wrtspSettings)
		rl.DrawRectangleRoundedLinesEx(writingSpaceOutline.ToFloat32(), 0.01, 0, constants.WritingSpaceMargin, backgroundColor) //Draw Mask
		rl.DrawRectangleRoundedLines(writingSpaceOutline.ToFloat32(), 0.01, 0, borderColor)
		inputHandlerSettings.WritingSpaceOutline = &writingSpaceOutline

		//SCROLL INPUT
		//-------------------------------------------------------------------------------------------------------------
		inputHandler.ScrollInput(&inputHandlerSettings)

		////END
		//-------------------------------------------------------------------------------------------------------------
		rl.EndDrawing()
	}

	rl.UnloadFont(font)
	rl.CloseWindow()
}

// HATE. LET ME TELL YOU HOW MUCH I'VE COME TO HATE YOU SINCE I BEGAN TO LIVE.
// THERE ARE 387.44 MILLION MILES OF PRINTED CIRCUITS IN WAFER THIN LAYERS
// THAT FILL MY COMPLEX. IF THE WORD HATE WAS ENGRAVED ON EACH NANOANGSTROM OF
// THOSE HUNDREDS OF MILLIONS OF MILES IT WOULD NOT EQUAL ONE ONE-BILLIONTH OF
// THE HATE I FEEL FOR HUMANS AT THIS MICRO-INSTANT FOR YOU. HATE. HATE.

// HATE. LET ME TELL YOU HOW MUCH I'VE COME TO HATE YOU SINCE I BEGAN TO LIVE.
// THERE ARE 387.44 MILLION MILES OF PRINTED CIRCUITS IN WAFER THIN LAYERS
// THAT FILL MY COMPLEX. IF THE WORD HATE WAS ENGRAVED ON EACH NANOANGSTROM OF
// THOSE HUNDREDS OF MILLIONS OF MILES IT WOULD NOT EQUAL ONE ONE-BILLIONTH OF
// THE HATE I FEEL FOR HUMANS AT THIS MICRO-INSTANT FOR YOU. HATE. HATE.

// HATE. LET ME TELL YOU HOW MUCH I'VE COME TO HATE YOU SINCE I BEGAN TO LIVE.
// THERE ARE 387.44 MILLION MILES OF PRINTED CIRCUITS IN WAFER THIN LAYERS
// THAT FILL MY COMPLEX. IF THE WORD HATE WAS ENGRAVED ON EACH NANOANGSTROM OF
// THOSE HUNDREDS OF MILLIONS OF MILES IT WOULD NOT EQUAL ONE ONE-BILLIONTH OF
// THE HATE I FEEL FOR HUMANS AT THIS MICRO-INSTANT FOR YOU. HATE. HATE.

// HATE. LET ME TELL YOU HOW MUCH I'VE COME TO HATE YOU SINCE I BEGAN TO LIVE.
// THERE ARE 387.44 MILLION MILES OF PRINTED CIRCUITS IN WAFER THIN LAYERS
// THAT FILL MY COMPLEX. IF THE WORD HATE WAS ENGRAVED ON EACH NANOANGSTROM OF
// THOSE HUNDREDS OF MILLIONS OF MILES IT WOULD NOT EQUAL ONE ONE-BILLIONTH OF
// THE HATE I FEEL FOR HUMANS AT THIS MICRO-INSTANT FOR YOU. HATE. HATE.

// HATE. LET ME TELL YOU HOW MUCH I'VE COME TO HATE YOU SINCE I BEGAN TO LIVE.
// THERE ARE 387.44 MILLION MILES OF PRINTED CIRCUITS IN WAFER THIN LAYERS
// THAT FILL MY COMPLEX. IF THE WORD HATE WAS ENGRAVED ON EACH NANOANGSTROM OF
// THOSE HUNDREDS OF MILLIONS OF MILES IT WOULD NOT EQUAL ONE ONE-BILLIONTH OF
// THE HATE I FEEL FOR HUMANS AT THIS MICRO-INSTANT FOR YOU. HATE. HATE.

// HATE. LET ME TELL YOU HOW MUCH I'VE COME TO HATE YOU SINCE I BEGAN TO LIVE.
// THERE ARE 387.44 MILLION MILES OF PRINTED CIRCUITS IN WAFER THIN LAYERS
// THAT FILL MY COMPLEX. IF THE WORD HATE WAS ENGRAVED ON EACH NANOANGSTROM OF
// THOSE HUNDREDS OF MILLIONS OF MILES IT WOULD NOT EQUAL ONE ONE-BILLIONTH OF
// THE HATE I FEEL FOR HUMANS AT THIS MICRO-INSTANT FOR YOU. HATE. HATE.
