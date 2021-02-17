package main

/*
TODO:
-Ambient game music / Fix ghost music
-Endings
-Multiply damage/speed in ascended mode
-Pause screen
	-Restart game
	-Mute music/sfx
	-Controls list
*/

//Average playthrough ~25 mins

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
	rand.Seed(time.Now().UnixNano())

	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowResizable(true)
	ebiten.SetWindowTitle("Feta Feles Remake")
	ebiten.SetRunnableOnUnfocused(true)

	ChangeAppState(new(TitleScreen))

	if err := ebiten.RunGame(new(App)); err != nil {
		log.Fatal(err)
	}
}
