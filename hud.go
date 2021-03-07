package main

import (
	"fmt"
	"image"

	"github.com/hajimehoshi/ebiten"
)

type GameHUD struct {
	root      *UINode
	loveBar   *UIBox
	msgText   *UIText
	msgTimer  float64
	timerText *UIText
	fpsText   *UIText
	pause     PauseScreen
}

type PauseScreen struct {
	container    *UIBox
	controlsButt *UIBox
	restartButt  *UIBox
}

func CreateGameHUD() *GameHUD {
	hud := &GameHUD{}
	hud.root = EmptyUINode()

	loveBorder := CreateUIBox(image.Rect(64, 40, 88, 48), image.Rect(4, 4, 4+160, 4+16), true)
	hud.root.AddChild(&loveBorder.UINode)
	hud.loveBar = CreateUIBox(image.Rect(104, 40, 112, 48), image.Rect(4, 4, loveBorder.Width()-4, loveBorder.Height()-4), false)
	loveBorder.AddChild(&hud.loveBar.UINode)

	msgBorder := CreateUIBox(image.Rect(112, 40, 136, 48), image.Rect(SCR_WIDTH_H-88, SCR_HEIGHT-48, SCR_WIDTH_H+88, SCR_HEIGHT-16), true)
	msgBorder.visible = false
	hud.root.AddChild(&msgBorder.UINode)
	hud.msgText = GenerateText("HEEBY DEEBY", image.Rect(8, 8, msgBorder.Width()-8, msgBorder.Height()-8))
	msgBorder.AddChild(&hud.msgText.UINode)

	timerBorder := CreateUIBox(image.Rect(112, 40, 136, 48), image.Rect(4.0, 20.0, 100.0, 36.0), true)
	hud.root.AddChild(&timerBorder.UINode)
	hud.timerText = GenerateText("00:00/00:00", image.Rect(4, 4, 2048, 2048))
	timerBorder.AddChild(&hud.timerText.UINode)

	hud.fpsText = GenerateText("FPS: 00", image.Rect(SCR_WIDTH-80, 0, SCR_WIDTH, 64))
	hud.root.AddChild(&hud.fpsText.UINode)

	hud.pause.container = CreateUIBox(image.Rect(136, 40, 160, 48), image.Rect(80, 48, 240, 176), true) //Background panel
	hud.pause.container.visible = false

	titleBox := CreateUIBox(image.Rect(112, 40, 136, 48), image.Rect(0, 0, 64, 16), true) //Header
	titleBox.AddChild(&GenerateText("PAUSE", image.Rect(8, 4, 2048, 2048)).UINode)
	hud.pause.container.AddChild(&titleBox.UINode)

	hud.pause.controlsButt = CreateUIBox(image.Rect(88, 40, 112, 48), image.Rect(0, 0, 108, 16), true) //Controls button
	hud.pause.controlsButt.AddChild(&GenerateText("CONTROLS", image.Rect(4, 4, 2048, 2048)).UINode)
	hud.pause.container.AddChild(&hud.pause.controlsButt.UINode)

	hud.pause.restartButt = CreateUIBox(image.Rect(88, 40, 112, 48), image.Rect(0, 0, 108, 16), true) //Restart button
	hud.pause.restartButt.AddChild(&GenerateText("RESTART GAME", image.Rect(4, 4, 2048, 2048)).UINode)
	hud.pause.container.AddChild(&hud.pause.restartButt.UINode)

	hud.pause.container.ArrangeChildren(image.Rect(4, 16, 4, 16))

	hud.root.AddChild(&hud.pause.container.UINode)

	return hud
}

func (hud *GameHUD) Update(game *Game) {
	if game.pause {
		hud.pause.container.visible = true
		//Respond to pause screen buttons
		if hud.pause.restartButt.Clicked() {
			ChangeAppState(NewGame(0))
		}
	} else {
		hud.pause.container.visible = false

		//Update UI with love amount
		barRect := image.Rect(3, 3, hud.loveBar.parent.Width()-3, hud.loveBar.parent.Height()-3) //This is the maximum size
		quotaPercent := float64(game.love) / float64(game.mission.loveQuota)
		barRect.Max.X = barRect.Min.X + int(float64(barRect.Size().X)*quotaPercent) //Resize based on amount of love
		hud.loveBar.dest = barRect
		hud.loveBar.Regen()

		//Update gameplay timer
		tSeconds := int(game.elapsedTime) % 60
		tMinutes := int(game.elapsedTime / 60.0)
		pSeconds := game.mission.parTime % 60
		pMinutes := game.mission.parTime / 60
		hud.timerText.text = fmt.Sprintf("%02d:%02d/%02d:%02d", tMinutes, tSeconds, pMinutes, pSeconds)
		hud.timerText.Regen()

		//Update message timer
		if hud.msgText != nil && hud.msgTimer > 0.0 {
			hud.msgText.parent.visible = true
			hud.msgText.Update(game.deltaTime)
			hud.msgTimer -= game.deltaTime
		} else {
			hud.msgText.parent.visible = false
		}
	}

	//FPS counter
	if debugDraw {
		hud.fpsText.visible = true
		hud.fpsText.text = fmt.Sprintf("FPS: %.2f", ebiten.CurrentTPS())
		hud.fpsText.Regen()
	} else {
		hud.fpsText.visible = false
	}
}

func (hud *GameHUD) Draw(screen *ebiten.Image) {
	hud.root.Draw(screen, nil)
}

func (hud *GameHUD) DisplayMessage(msg string, time float64) {
	hud.msgTimer = time
	hud.msgText.text = msg
	hud.msgText.Regen()
	hud.msgText.fillPos = 0
}
