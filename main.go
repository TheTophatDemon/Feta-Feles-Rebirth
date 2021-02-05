package main

/*
TODO:
-Add failsafe for when cat gets stuck in a wall
-Enemies still getting stuck in eachother
-Loud explosions bug
-Secret teleporter / Powerup
-Ending screen
-Some sort of end-game skill assessment
-Worm
-Music
*/

import (
	"image"
	_ "image/color"
	_ "image/png"
	"log"
	"math"
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

func init() {
	__lastTime = time.Now()
}

const FRAMERATE = 1.0 / 60.0

func (a *App) Update() error {
	now := time.Now()
	deltaTime := now.Sub(__lastTime).Seconds()
	__lastTime = now

	__appState.Update(math.Min(FRAMERATE, deltaTime))

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
	ebiten.SetWindowTitle("Feta Feles Remake")
	ebiten.SetRunnableOnUnfocused(true)

	ChangeAppState(new(TitleScreen))

	if err := ebiten.RunGame(new(App)); err != nil {
		log.Fatal(err)
	}
}

var whiteFadeShader *ebiten.Shader
var noiseImg *ebiten.Image

func init() {
	noiseImg = ebiten.NewImage(SCR_WIDTH, SCR_HEIGHT)
	noisePixels := make([]byte, SCR_WIDTH*SCR_HEIGHT*4)
	for i := 0; i < SCR_HEIGHT*SCR_WIDTH; i++ {
		noisePixels[i*4+0] = byte(rand.Intn(255))
		noisePixels[i*4+1] = byte(rand.Intn(255))
		noisePixels[i*4+2] = byte(rand.Intn(255))
		noisePixels[i*4+3] = 255
	}
	noiseImg.ReplacePixels(noisePixels)

	var err error
	whiteFadeShader, err = ebiten.NewShader([]byte(`
		package main

		var Coverage float

		func Fragment(position vec4, texCoord vec2, color vec4) vec4 {
			diffuse := imageSrc0UnsafeAt(texCoord)
			noise := imageSrc1UnsafeAt(texCoord)
			mask := step(1.0 - Coverage, noise.r)
			return min(diffuse + mask, vec4(0.9, 0.9, 0.9, 1.0))
		}
	`))
	if err != nil {
		println(err)
		panic(err)
	}
}
