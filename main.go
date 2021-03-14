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
	"image"
	_ "image/color"
	_ "image/png"
	"log"
	"math"
	"math/rand"

	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/thetophatdemon/Feta-Feles-Remastered/assets"
	"github.com/thetophatdemon/Feta-Feles-Remastered/audio"
)

const (
	SCR_WIDTH    = 320
	SCR_WIDTH_H  = SCR_WIDTH / 2
	SCR_HEIGHT   = 240
	SCR_HEIGHT_H = SCR_HEIGHT / 2
)

type AppState interface {
	Update(deltaTime float64)
	Draw(screen *ebiten.Image)
	Enter()
	Leave()
}

type App struct{}

func (a *App) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return SCR_WIDTH, SCR_HEIGHT
}

var __lastTime time.Time
var __appState AppState

func init() {
	__lastTime = time.Now()
}

const FRAMERATE = 1.0 / 60.0

func (a *App) Update() error {
	now := time.Now()
	deltaTime := now.Sub(__lastTime).Seconds()
	__lastTime = now

	__appState.Update(math.Min(FRAMERATE, deltaTime))
	audio.Update(deltaTime)

	//Toggle fullscreen with alt + enter
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) && ebiten.IsKeyPressed(ebiten.KeyAlt) {
		ebiten.SetFullscreen(!ebiten.IsFullscreen())
	}

	return nil
}

//Draw ...
func (a *App) Draw(screen *ebiten.Image) {
	__appState.Draw(screen)
}

func ChangeAppState(newState AppState) {
	if newState == nil {
		panic("Can't change into nil app state!")
	}
	if __appState != nil {
		__appState.Leave()
	}
	__appState = newState
	newState.Enter()
}

var __graphics *ebiten.Image

//Returns the graphics page and loads it if it isn't there
func GetGraphics() *ebiten.Image {
	if __graphics == nil {
		img, _, err := image.Decode(assets.ReadCompressedString(assets.PNG_GRAPHICS))
		if err != nil {
			log.Fatal(err)
		}
		__graphics = ebiten.NewImageFromImage(img)
	}
	return __graphics
}

func main() {
	seed := time.Now().UnixNano() % 1615698000000000000
	rand.Seed(seed)

	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowResizable(true)
	ebiten.SetWindowTitle("Feta Feles Rebirth")
	ebiten.SetRunnableOnUnfocused(true)

	ChangeAppState(new(TitleScreen))

	if err := ebiten.RunGame(new(App)); err != nil {
		log.Fatal(err)
	}
}
