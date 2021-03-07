package main

import (
	"image"

	"github.com/hajimehoshi/ebiten"
)

type PauseScreen struct {
	background  *UIBox
	titleBox    *UIBox
	controlsBox *UIBox
	restartBox  *UIBox
}

func NewPauseScreen() *PauseScreen {
	background := CreateUIBox(image.Rect(136, 40, 160, 48), image.Rect(80, 48, 240, 176))
	background.padding = image.Rect(4, 16, 4, 16)
	titleBox := CreateUIBox(image.Rect(112, 40, 136, 48), image.Rect(0, 0, 64, 16))
	titleBox.children.PushBack(GenerateText("PAUSE", image.Rect(8, 4, titleBox.size.X-8, titleBox.size.Y-4)))
	background.children.PushBack(titleBox)

	controlsBox := CreateUIBox(image.Rect(88, 40, 112, 48), image.Rect(0, 0, 108, 16))
	controlsBox.children.PushBack(GenerateText("CONTROLS", image.Rect(4, 4, controlsBox.size.X-4, controlsBox.size.Y-4)))
	background.children.PushBack(controlsBox)

	restartBox := CreateUIBox(image.Rect(88, 40, 112, 48), image.Rect(0, 0, 108, 16))
	restartBox.children.PushBack(GenerateText("RESTART GAME", image.Rect(4, 4, restartBox.size.X-4, restartBox.size.Y-4)))
	background.children.PushBack(restartBox)

	background.ArrangeChildren()

	return &PauseScreen{
		background:  background,
		titleBox:    titleBox,
		controlsBox: controlsBox,
		restartBox:  restartBox,
	}
}

func (ps *PauseScreen) Update(deltaTime float64) {
	switch {
	case ps.restartBox.IsClicked():
		ChangeAppState(NewGame(0))
	}
}

func (ps *PauseScreen) Draw(screen *ebiten.Image) {
	ps.background.Draw(screen, nil)
}
