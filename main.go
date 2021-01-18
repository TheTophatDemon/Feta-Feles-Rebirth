package main

/*
TODO:
-Powerups in little caverns?
-Worm
-Bullet fx
-Title Screen
-Mission 0: Show basics of game
-Feles
-Music
-Store assets as embedded zip file...?
-Fix jitter at beginning
-Fix cats getting stuck in walls...?
*/

import (
	"bytes"
	"image"
	_ "image/color"
	_ "image/png"
	"log"
	"math/rand"

	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/thetophatdemon/Feta-Feles-Remastered/assets"
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

func (a *App) Update() error {
	now := time.Now()
	deltaTime := now.Sub(__lastTime).Seconds()
	__lastTime = now

	__appState.Update(deltaTime)

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
		img, _, err := image.Decode(bytes.NewReader(assets.Parse(assets.PNG_GRAPHICS)))
		if err != nil {
			log.Fatal(err)
		}
		__graphics = ebiten.NewImageFromImage(img)
	}
	return __graphics
}

func main() {
	rand.Seed(time.Now().UnixNano())

	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Feta Feles Remake")
	ebiten.SetRunnableOnUnfocused(true)

	ChangeAppState(NewGame(0))

	if err := ebiten.RunGame(new(App)); err != nil {
		log.Fatal(err)
	}
}
