/*
Copyright (C) 2021 Alexander Lunsford

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package main

import (
	"fmt"
	"image"

	"github.com/hajimehoshi/ebiten"
	"github.com/thetophatdemon/Feta-Feles-Remastered/audio"
)

type GameHUD struct {
	root      *UINode
	loveBar   *UIBox
	msgText   *UIText
	msgTimer  float64
	timerText *UIText
	fpsText   *UIText
	menu      *UINode
	pause     PauseScreen
	control   ControlsScreen
}

type PauseScreen struct {
	container    *UIBox
	controlsButt *UIBox
	restartButt  *UIBox
	musicButt    *UIBox
	sfxButt      *UIBox
}

type ControlsScreen struct {
	container *UIBox
	backButt  *UIBox
}

func CreateGameHUD() *GameHUD {
	hud := &GameHUD{}
	hud.root = EmptyUINode()

	//==============================
	//IN GAME
	//==============================

	loveBorder := CreateUIBox(image.Rect(64, 40, 88, 48), image.Rect(4, 4, 4+160, 4+16), true)
	hud.root.AddChild(&loveBorder.UINode)
	hud.loveBar = CreateUIBox(image.Rect(104, 40, 112, 48), image.Rect(4, 4, loveBorder.Width()-4, loveBorder.Height()-4), false)
	loveBorder.AddChild(&hud.loveBar.UINode)

	msgBorder := CreateUIBox(image.Rect(112, 40, 136, 48), image.Rect(SCR_WIDTH_H-88, SCR_HEIGHT-48, SCR_WIDTH_H+88, SCR_HEIGHT-8), true)
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

	//Add background panel for menus
	hud.menu = EmptyUINode()
	hud.menu.visible = false
	hud.root.AddChild(hud.menu)

	//=====================================
	//PAUSE SCREEN
	//=====================================

	hud.pause.container = CreateUIBox(image.Rect(136, 40, 160, 48), image.Rect(80, 48, 240, 176), true)
	hud.pause.container.visible = true //This is overwritten by the menu's visibility, but important for keeping track of menu state

	titleBox := CreateUIBox(image.Rect(112, 40, 136, 48), image.Rect(0, 0, 64, 16), true) //Header
	titleBox.AddChild(&GenerateText("PAUSE", image.Rect(8, 4, 2048, 2048)).UINode)
	hud.pause.container.AddChild(&titleBox.UINode)

	muteContainer := CreateUIBox(image.Rect(152, 40, 160, 48), image.Rect(0, 0, 128, 16), false)
	hud.pause.container.AddChild(&muteContainer.UINode)

	hud.pause.sfxButt = CreateUIBox(image.Rect(128, 160, 144, 176), image.Rect(0, 0, 16, 16), false)
	muteContainer.AddChild(&hud.pause.sfxButt.UINode)
	cross := CreateUIBox(image.Rect(160, 160, 176, 176), image.Rect(0, 0, 16, 16), false)
	cross.visible = audio.MuteSfx
	hud.pause.sfxButt.AddChild(&cross.UINode)

	hud.pause.musicButt = CreateUIBox(image.Rect(144, 160, 160, 176), image.Rect(0, 0, 16, 16), false)
	muteContainer.AddChild(&hud.pause.musicButt.UINode)
	cross = CreateUIBox(image.Rect(160, 160, 176, 176), image.Rect(0, 0, 16, 16), false)
	cross.visible = audio.MuteMusic
	hud.pause.musicButt.AddChild(&cross.UINode)

	muteContainer.ArrangeChildren(image.Rect(0, 0, 0, 0), false)

	hud.pause.controlsButt = CreateUIBox(image.Rect(88, 40, 112, 48), image.Rect(0, 0, 108, 16), true) //Controls button
	hud.pause.controlsButt.AddChild(&GenerateText("HELP", image.Rect(4, 4, 2048, 2048)).UINode)
	hud.pause.container.AddChild(&hud.pause.controlsButt.UINode)

	hud.pause.restartButt = CreateUIBox(image.Rect(88, 40, 112, 48), image.Rect(0, 0, 108, 16), true) //Restart button
	hud.pause.restartButt.AddChild(&GenerateText("RESTART GAME", image.Rect(4, 4, 2048, 2048)).UINode)
	hud.pause.container.AddChild(&hud.pause.restartButt.UINode)

	hud.pause.container.ArrangeChildren(image.Rect(4, 4, 4, 8), true)

	hud.menu.AddChild(&hud.pause.container.UINode)

	//================================
	//CONTROLS SCREEN
	//===============================

	hud.control.container = CreateUIBox(image.Rect(136, 40, 160, 48), image.Rect(16, 40, SCR_WIDTH-16, SCR_HEIGHT-8), true)
	hud.control.container.visible = false

	titleBox = CreateUIBox(image.Rect(112, 40, 136, 48), image.Rect(0, 0, 80, 16), true) //Header
	titleBox.AddChild(&GenerateText("TIPS", image.Rect(8, 4, 2048, 2048)).UINode)
	hud.control.container.AddChild(&titleBox.UINode)

	lineRect := image.Rect(0, 0, SCR_WIDTH-32-16, 20)
	controlsText := []*UIText{
		GenerateText("SPACE KEY WILL FIX SHOOTING       DIRECTION", lineRect),
		GenerateText("CAT IS ONLY HURT BY ASCENDED      BULLETS", lineRect),
		GenerateText("LISTEN FOR ITS MEWS", lineRect),
		GenerateText("PUSH BOUNDARIES TO WARP", lineRect),
		GenerateText("BOUNCY BULLETS EXPLODE RUNES", lineRect),
		GenerateText("EVERY MISSION HAS A PAR TIME", lineRect),
		GenerateText("BEAT THE PAR TIMES FOR EVERY      MISSION IF YOU DARE", lineRect),
	}
	for _, t := range controlsText {
		hud.control.container.AddChild(&t.UINode)
	}

	hud.control.backButt = CreateUIBox(image.Rect(88, 40, 112, 48), image.Rect(0, 0, 108, 16), true) //Back button
	hud.control.backButt.AddChild(&GenerateText("RETURN", image.Rect(4, 4, 2048, 2048)).UINode)
	hud.control.container.AddChild(&hud.control.backButt.UINode)

	hud.control.container.ArrangeChildren(image.Rect(4, 4, 4, 8), true)

	hud.menu.AddChild(&hud.control.container.UINode)

	return hud
}

func (hud *GameHUD) Update(game *Game) {
	if game.pause {
		hud.menu.visible = true
		if hud.pause.container.visible {
			//Respond to pause screen buttons
			if hud.pause.restartButt.Clicked() {
				audio.PlaySound("button")
				ChangeAppState(NewGame(0))
			} else if hud.pause.sfxButt.Clicked() {
				audio.PlaySound("button")
				audio.MuteSfx = !audio.MuteSfx
				check := hud.pause.sfxButt.children.Front().Value.(*UINode)
				check.visible = audio.MuteSfx
			} else if hud.pause.musicButt.Clicked() {
				audio.PlaySound("button")
				audio.MuteMusic = !audio.MuteMusic
				check := hud.pause.musicButt.children.Front().Value.(*UINode)
				check.visible = audio.MuteMusic
			} else if hud.pause.controlsButt.Clicked() {
				audio.PlaySound("button")
				hud.pause.container.visible = false
				hud.control.container.visible = true
			}
		} else if hud.control.container.visible {
			//Respond to controls screen buttons
			if hud.control.backButt.Clicked() {
				audio.PlaySound("button")
				hud.control.container.visible = false
				hud.pause.container.visible = true
			}
		}
	} else {
		hud.menu.visible = false

		//Hide love bar if player is under it
		hud.loveBar.parent.visible = (game.playerObj.pos.X > 164 || game.playerObj.pos.Y > 64)
		hud.timerText.parent.visible = hud.loveBar.parent.visible

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

func (hud *GameHUD) IsDisplayingMessage() bool {
	return hud.msgTimer > 0.0
}